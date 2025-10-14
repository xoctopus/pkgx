package pkgx

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"iter"
	"maps"
	"path/filepath"
	"slices"
	"sort"

	"github.com/xoctopus/x/mapx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/ptrx"
	gopkg "golang.org/x/tools/go/packages"

	"github.com/xoctopus/pkgx/internal"
)

type (
	Doc       = internal.Doc
	ModuleSum = internal.Sum
	Constant  = internal.Constant
	Function  = internal.Function
	TypeName  = internal.TypeName

	Constants = internal.Objects[*types.Const, *Constant]
	Functions = internal.Objects[*types.Func, *Function]
	TypeNames = internal.Objects[*types.TypeName, *TypeName]

	TPackage  = types.Package
	GoPackage = gopkg.Package
	GoModule  = gopkg.Module
)

func NewPackages(patterns ...string) *Packages {
	u := &Packages{
		entries:  patterns,
		fileset:  token.NewFileSet(),
		packages: mapx.NewSafeXmap[string, Package](),
		modules:  mapx.NewSafeSet[string](),
		directs:  mapx.NewSafeSet[string](),
		sums:     mapx.NewSafeXmap[string, ModuleSum](),
	}

	packages, err := gopkg.Load(&gopkg.Config{
		Fset: u.fileset,
		Mode: gopkg.LoadMode(0b11111111111111111),
	}, patterns...)
	must.NoErrorF(err, "failed to load packages: %v", patterns)

	var register func(p *GoPackage)

	register = func(p *GoPackage) {
		x := newx(p)
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
				s, _ := u.sums.LoadOrStore(p.Module.Path, internal.NewSum(p.Module.Dir))
				s.Add(p)
			}
		}
	}

	sort.Slice(packages, func(i, j int) bool {
		return packages[i].PkgPath < packages[j].PkgPath
	})

	for _, p := range packages {
		must.BeTrueF(len(p.Errors) == 0, "loaded package `%s` error", p.PkgPath)
		if p.Module != nil {
			u.modules.Store(p.Module.Path)
		}
		u.directs.Store(p.PkgPath)
	}

	for _, p := range packages {
		register(p)
	}

	return u
}

type Packages struct {
	entries  []string
	fileset  *token.FileSet
	packages mapx.Map[string, Package]
	modules  mapx.Set[string]
	directs  mapx.Set[string]
	sums     mapx.Map[string, ModuleSum]
}

// Package locates package by path
func (u *Packages) Package(path string) Package {
	p, _ := u.packages.Load(path)
	return p
}

// ModuleSum returns module sum by module path
func (u *Packages) ModuleSum(module string) ModuleSum {
	s, _ := u.sums.Load(module)
	return s
}

// Packages returns iteration for all loaded package, include package from std and general importing
func (u *Packages) Packages() iter.Seq2[string, Package] {
	return u.packages.Range
}

// Directs returns iteration for packages under entries
func (u *Packages) Directs() iter.Seq[string] {
	return u.directs.Range
}

// Modules returns iteration for modules under entries
func (u *Packages) Modules() iter.Seq[string] {
	return u.modules.Range
}

type Package interface {
	Unwrap() *TPackage
	GoPackage() *GoPackage
	GoModule() *GoModule
	PackageByPath(string) Package
	Docs() iter.Seq[*Doc]

	SourceDir() string
	Eval(ast.Expr) (types.TypeAndValue, error)
	Files() []*ast.File
	FileSet() *token.FileSet
	Position(token.Pos) token.Position
	ObjectOf(*ast.Ident) types.Object

	TypeNames() TypeNames
	Constants() Constants
	Functions() Functions
}

func newx(p *gopkg.Package) Package {
	must.BeTrue(p != nil && len(p.Errors) == 0)
	x := &xpkg{
		p:         p,
		imports:   mapx.NewXmap[string, Package](),
		typenames: internal.NewObjects[*types.TypeName, *TypeName](),
		constants: internal.NewObjects[*types.Const, *Constant](),
		functions: internal.NewObjects[*types.Func, *Function](),
	}
	methods := make(map[types.Type][]*Function)

	for _, file := range p.Syntax {
		if file.Doc != nil {
			x.docs = append(x.docs, internal.ParseDocument(file.Doc))
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch d := node.(type) {
			case *ast.GenDecl:
				if d.Tok != token.TYPE && d.Tok != token.CONST {
					return true
				}
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ValueSpec:
						doc := internal.ParseDocument(d.Doc, s.Doc, s.Comment)
						for _, ident := range s.Names {
							x.constants.Add(&Constant{
								Object: internal.NewObject(
									s,
									ident,
									p.TypesInfo.Defs[ident].(*types.Const),
									doc,
								),
							})
						}
					case *ast.TypeSpec:
						x.typenames.Add(internal.NewTypeName(
							internal.NewObject(
								s,
								s.Name,
								p.TypesInfo.Defs[s.Name].(*types.TypeName),
								internal.ParseDocument(d.Doc, s.Doc, s.Comment),
							),
						))
					}
				}
			case *ast.FuncDecl:
				doc := internal.ParseDocument(d.Doc)
				u := p.TypesInfo.Defs[d.Name].(*types.Func)
				o := internal.NewObject(node, d.Name, u, doc)
				f := &internal.Function{Object: o}
				if recv := u.Signature().Recv(); recv == nil {
					x.functions.Add(f)
				} else {
					t := types.Unalias(internal.Deref(recv.Type()))
					methods[t] = append(methods[t], f)
				}
			}
			return true
		})
	}

	x.functions.Init()
	x.constants.Init()
	x.typenames.Init()

	for t := range x.typenames.Elements() {
		t.AddMethods(methods[t.Type()]...)
	}
	sort.Slice(x.docs, func(i, j int) bool {
		return x.docs[i].Pos() < x.docs[j].Pos()
	})

	// TODO inspecting signatures should contains FuncDecl, FuncLit and CallExpr
	// TODO should analyze signatures returned results
	return x
}

type xpkg struct {
	p    *gopkg.Package
	u    *Packages
	dir  *string
	docs []*Doc

	// fileset *token.FileSet
	imports mapx.Map[string, Package]

	typenames TypeNames
	constants Constants
	functions Functions
	// TODO signatures and results
	// signatures internal.Objects[*types.Signature, *internal.Signature]
}

func (x *xpkg) Unwrap() *types.Package {
	return x.p.Types
}

func (x *xpkg) GoPackage() *gopkg.Package {
	return x.p
}

func (x *xpkg) GoModule() *gopkg.Module {
	return x.p.Module
}

func (x *xpkg) Docs() iter.Seq[*Doc] {
	return slices.Values(x.docs)
}

func (x *xpkg) PackageByPath(path string) Package {
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

func (x *xpkg) Constants() Constants {
	return x.constants
}

func (x *xpkg) Functions() Functions {
	return x.functions
}

func (x *xpkg) TypeNames() TypeNames {
	return x.typenames
}
