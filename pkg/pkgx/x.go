package pkgx

import (
	"bytes"
	"context"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"maps"
	"path/filepath"
	"slices"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/ptrx"
	"github.com/xoctopus/x/syncx"
	gopkg "golang.org/x/tools/go/packages"

	internal "github.com/xoctopus/pkgx/internal/pkgx"
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

func NewPackages(ctx context.Context, patterns ...string) *Packages {
	u := &Packages{
		entries:  patterns,
		fileset:  token.NewFileSet(),
		packages: syncx.NewXmap[string, Package](),
		modules:  syncx.NewSet[string](),
		directs:  syncx.NewSet[string](),
		sums:     syncx.NewXmap[string, ModuleSum](),
	}
	ctx = CtxFileset.With(ctx, u.fileset)

	packages, err := gopkg.Load(Config(ctx), patterns...)
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

	for _, p := range packages {
		must.BeTrueF(len(p.Errors) == 0, "loaded package `%s` error: %v", p.ID, p.Errors)
		if p.Module != nil {
			u.modules.Store(p.Module.Path)
		}
		u.directs.Store(p.PkgPath)
	}

	for _, p := range packages {
		register(p)
	}

	for _, p := range u.Packages {
		x := p.(*xpkg)
		x.typenames.(internal.ObjectsManager[*types.TypeName, *TypeName]).Init(u.fileset)
		x.functions.(internal.ObjectsManager[*types.Func, *Function]).Init(u.fileset)
		x.constants.(internal.ObjectsManager[*types.Const, *Constant]).Init(u.fileset)
	}

	return u
}

type Packages struct {
	entries  []string
	fileset  *token.FileSet
	packages syncx.Map[string, Package]
	modules  *syncx.Set[string]
	directs  *syncx.Set[string]
	sums     syncx.Map[string, ModuleSum]
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
func (u *Packages) Packages(f func(string, Package) bool) {
	u.packages.Range(f)
}

// Directs returns iteration for packages under entries
func (u *Packages) Directs(f func(string) bool) {
	u.directs.Range(f)
}

// Modules returns iteration for modules under entries
func (u *Packages) Modules(f func(string) bool) {
	u.modules.Range(f)
}

func (u *Packages) DocOf(pos token.Pos) *Doc {
	for _, p := range u.Packages {
		if d := p.DocOf(pos); d != nil {
			return d
		}
	}
	return nil
}

type Package interface {
	Path() string
	ID() string
	Name() string

	// Unwrap returns types.Package of this package
	Unwrap() *TPackage
	// GoPackage returns packages.Package of this package
	GoPackage() *GoPackage
	// GoModule returns packages.Module of this package
	GoModule() *GoModule
	// PackageByPath locates package of given path
	PackageByPath(string) Package
	// Doc returns package level documents
	Doc() *Doc
	// SourceDir returns dir path of current package
	SourceDir() string
	// DocOf returns doc of node
	DocOf(token.Pos) *Doc

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
		docs:      syncx.NewXmap[token.Pos, *Doc](),
		imports:   syncx.NewXmap[string, Package](),
		typenames: internal.NewObjects[*types.TypeName, *TypeName](),
		constants: internal.NewObjects[*types.Const, *Constant](),
		functions: internal.NewObjects[*types.Func, *Function](),
	}
	methods := make(map[types.Type][]*Function)
	docs := make([]*ast.CommentGroup, len(p.Syntax))

	for _, file := range p.Syntax {
		if file.Doc != nil {
			docs = append(docs, file.Doc)
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
						x.docs.Store(s.Pos(), doc)
						for _, ident := range s.Names {
							x.constants.(internal.ObjectsManager[*types.Const, *Constant]).
								Add(&Constant{
									Object: internal.NewObject(
										s,
										ident,
										p.TypesInfo.Defs[ident].(*types.Const),
										doc,
									),
								})
							x.docs.Store(ident.Pos(), doc)
						}
					case *ast.TypeSpec:
						doc := internal.ParseDocument(d.Doc, s.Doc, s.Comment)
						x.docs.Store(s.Pos(), doc)
						x.typenames.(internal.ObjectsManager[*types.TypeName, *TypeName]).
							Add(internal.NewTypeName(
								internal.NewObject(
									s,
									s.Name,
									p.TypesInfo.Defs[s.Name].(*types.TypeName),
									doc,
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
					x.functions.(internal.ObjectsManager[*types.Func, *Function]).Add(f)
				} else {
					t := types.Unalias(internal.Deref(recv.Type()))
					methods[t] = append(methods[t], f)
				}
			case *ast.StructType:
				for _, f := range d.Fields.List {
					doc := internal.ParseDocument(f.Doc, f.Comment)
					if len(f.Names) == 0 {
						pos := token.NoPos
						exp := f.Type
						for pos == token.NoPos {
							switch u := exp.(type) {
							case *ast.Ident:
								pos = u.Pos()
							case *ast.SelectorExpr:
								pos = u.Sel.Pos()
							case *ast.StarExpr:
								exp = u.X
							case *ast.IndexExpr:
								exp = u.X
							default:
								_, ok := u.(*ast.IndexListExpr)
								must.BeTrueF(ok, "unexpected ast type as struct field: %T", u)
								exp = u.(*ast.IndexListExpr).X
							}
						}
						x.docs.Store(pos, doc)
					} else {
						x.docs.Store(f.Pos(), doc)
					}
				}
			}
			return true
		})
	}

	if len(docs) > 0 {
		x.doc = internal.ParseDocument(docs[0], docs[1:]...)
	} else {
		x.doc = internal.DefaultDoc
	}

	typenames := x.typenames.(internal.ObjectsManager[*types.TypeName, *TypeName])
	for _, t := range typenames.RangeNodes {
		t.AddMethods(methods[t.Type()]...)
	}

	// TODO inspecting signatures should contains FuncDecl, FuncLit and CallExpr
	// TODO should analyze signatures returned results
	return x
}

type xpkg struct {
	p    *gopkg.Package
	u    *Packages
	dir  *string
	doc  *Doc
	docs syncx.Map[token.Pos, *Doc]

	// fileset *token.FileSet
	imports syncx.Map[string, Package]

	typenames TypeNames
	constants Constants
	functions Functions
	// TODO signatures and results
	// signatures internal.Objects[*types.Signature, *internal.Signature]
}

func (x *xpkg) Path() string {
	return x.p.PkgPath
}

func (x *xpkg) ID() string {
	return x.p.ID
}

func (x *xpkg) Name() string {
	return x.p.Name
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

func (x *xpkg) Doc() *Doc {
	return x.doc
}

func (x *xpkg) DocOf(pos token.Pos) *Doc {
	d, _ := x.docs.Load(pos)
	return d
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
