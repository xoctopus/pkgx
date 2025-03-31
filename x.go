package pkgx

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"maps"
	"path/filepath"
	"slices"
	"sort"
	"sync"

	"github.com/xoctopus/x/mapx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/ptrx"
	gopkg "golang.org/x/tools/go/packages"
)

func NewPackages(patterns ...string) *Packages {
	u := &Packages{
		fileset:  token.NewFileSet(),
		packages: mapx.NewXmap[string, Package](),
		modules:  mapx.NewSet[string](),
		directs:  mapx.NewSet[string](),
		sum:      &sum{hashes: make(map[string]string)},
	}

	loaded, err := gopkg.Load(&gopkg.Config{
		Fset: u.fileset,
		Mode: gopkg.LoadMode(0b11111111111111111),
	}, patterns...)
	must.NoErrorWrap(err, "failed to load packages: %v", patterns)

	var register func(p *gopkg.Package)

	register = func(p *gopkg.Package) {
		x := newX(p)
		x.(*xpkg).u = u

		for _, path := range slices.Sorted(maps.Keys(p.Imports)) {
			if _, ok := u.packages.Load(path); !ok {
				register(p.Imports[path])
			}
		}
		u.packages.Store(p.PkgPath, x)

		if p.Module != nil {
			if u.modules.Exists(p.Module.Path) {
				u.directs.Store(p.PkgPath)
				u.modules.Store(p.Module.Path)
				u.sum.AddPackage(p)
			}
		}
	}

	sort.Slice(loaded, func(i, j int) bool {
		return loaded[i].PkgPath < loaded[j].PkgPath
	})

	for _, p := range loaded {
		must.BeTrueWrap(len(p.Errors) == 0, "loaded package `%s` error", p.PkgPath)
		if p.Module != nil {
			u.modules.Store(p.Module.Path)
		}
		u.directs.Store(p.PkgPath)
	}

	for _, p := range loaded {
		register(p)
	}

	return u
}

type Packages struct {
	fileset  *token.FileSet
	packages mapx.Map[string, Package]
	modules  mapx.Set[string]
	directs  mapx.Set[string]
	sum      *sum
}

func (u *Packages) Package(path string) Package {
	p, _ := u.packages.Load(path)
	return p
}

func (u *Packages) Sum() Sum {
	return u.sum
}

type Package interface {
	Unwrap() *types.Package
	GoPackage() *gopkg.Package

	Package(string) Package
	Module() *gopkg.Module
	SourceDir() string
	Eval(ast.Expr) (types.TypeAndValue, error)
	Files() []*ast.File
	FileSet() *token.FileSet
	Position(token.Pos) token.Position
	ObjectOf(*ast.Ident) types.Object

	Doc(ast.Node) *Document

	Type(string) *Typename
	TypeByNode(ast.Node) *Typename
	NamedTypes() []*Typename

	Const(string) *Constant
	ConstByNode(ast.Node) *Constant
	Constants() []*Constant

	Func(string) *Function
	FuncByNode(n ast.Node) *Function
	Functions() []*Function
	Signatures() []*Signature
}

func newX(p *gopkg.Package) Package {
	must.BeTrue(p != nil && len(p.Errors) == 0)
	x := &xpkg{
		p: p,

		imports: mapx.NewXmap[string, Package](),

		documents:  ScanDocuments(p),
		typenames:  ScanTypenames(p),
		constants:  ScanConstants(p),
		functions:  ScanFunctions(p),
		signatures: ScanSignatures(p),
	}

	return x
}

type xpkg struct {
	p   *gopkg.Package
	u   *Packages
	dir *string

	fileset *token.FileSet
	imports mapx.Map[string, Package]

	documents *Set[*Document]
	typenames *Set[*Typename]
	constants *Set[*Constant]
	functions *Set[*Function]

	signatures *Set[*Signature]
	results    sync.Map
}

func (x *xpkg) Unwrap() *types.Package {
	return x.p.Types
}

func (x *xpkg) GoPackage() *gopkg.Package {
	return x.p
}

func (x *xpkg) Package(path string) Package {
	if _, ok := x.p.Imports[path]; !ok {
		return nil
	}

	if p, ok := x.imports.Load(path); ok {
		return p
	}

	p := x.u.Package(path)
	must.BeTrue(p != nil)
	x.imports.Store(path, p)
	return p
}

func (x *xpkg) SourceDir() (dir string) {
	if x.dir != nil {
		return *x.dir
	}

	defer func() {
		x.dir = ptrx.Ptr(dir)
	}()

	if x.p.Module == nil {
		return ""
	}
	if x.p.PkgPath == x.p.Module.Path {
		return x.p.Module.Dir
	}
	return filepath.Join(x.p.Module.Dir, x.p.PkgPath[len(x.p.Module.Path):])
}

func (x *xpkg) Module() *gopkg.Module {
	return x.p.Module
}

func (x *xpkg) Eval(e ast.Expr) (types.TypeAndValue, error) {
	code := bytes.NewBuffer(nil)
	if err := format.Node(code, x.p.Fset, e); err != nil {
		return types.TypeAndValue{}, err
	}

	return types.Eval(x.p.Fset, x.p.Types, e.Pos(), code.String())
}

func (x *xpkg) Files() []*ast.File {
	return x.p.Syntax
}

func (x *xpkg) FileSet() *token.FileSet {
	return x.u.fileset
}

func (x *xpkg) Position(p token.Pos) token.Position {
	return x.p.Fset.Position(p)
}

func (x *xpkg) ObjectOf(i *ast.Ident) types.Object {
	return x.p.TypesInfo.ObjectOf(i)
}

func (x *xpkg) Doc(n ast.Node) *Document {
	return x.documents.ValueByNode(n)
}

func (x *xpkg) Func(name string) *Function {
	return x.functions.ValueByName(name)
}

func (x *xpkg) FuncByNode(n ast.Node) *Function {
	return x.functions.ValueByNode(n)
}

func (x *xpkg) Functions() []*Function {
	return x.functions.Values()
}

func (x *xpkg) Signatures() []*Signature {
	return x.signatures.Values()
}

func (x *xpkg) Type(name string) *Typename {
	return x.typenames.ValueByName(name)
}

func (x *xpkg) TypeByNode(n ast.Node) *Typename {
	return x.typenames.ValueByNode(n)
}

func (x *xpkg) NamedTypes() []*Typename {
	return x.typenames.Values()
}

func (x *xpkg) Const(name string) *Constant {
	return x.constants.ValueByName(name)
}

func (x *xpkg) ConstByNode(n ast.Node) *Constant {
	return x.constants.ValueByNode(n)
}

func (x *xpkg) Constants() []*Constant {
	return x.constants.Values()
}
