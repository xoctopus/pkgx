package pkgx

import (
	"go/ast"
	"go/types"
	"sort"

	"github.com/xoctopus/x/mapx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
	gopkg "golang.org/x/tools/go/packages"
)

type Typename struct {
	node     ast.Node
	ident    *ast.Ident
	typename *types.TypeName
	typ      types.Type
	methods  mapx.Map[string, *Function]
}

func (t *Typename) IsZero() bool {
	return t == nil || t.typename == nil
}

func (t *Typename) Name() string {
	return t.ident.Name
}

func (t *Typename) Node() ast.Node {
	return t.node
}

func (t *Typename) Value() *types.TypeName {
	return t.typename
}

func (t *Typename) Type() types.Type {
	return t.typename.Type()
}

func (t *Typename) Methods() []*Function {
	methods := mapx.Values(t.methods)
	sort.Slice(methods, func(i, j int) bool {
		return methods[i].Name() < methods[j].Name()
	})
	return methods
}

func (t *Typename) MethodByName(name string) *Function {
	f, _ := t.methods.Load(name)
	return f
}

func ScanTypenames(pkg *gopkg.Package) *Set[*Typename] {
	typenames := NewSet[*Typename](pkg)

	for ident, obj := range pkg.TypesInfo.Defs {
		typename, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}
		if t := typename.Type(); !reflectx.CanCast[*types.Named](t) && !reflectx.CanCast[*types.Alias](t) {
			continue
		}
		node := IdentNode(ident)
		must.BeTrueWrap(node != nil, "%s: %s", ident.Name, SourceOfNode(pkg, node))
		elem := &Typename{
			node:     node,
			ident:    ident,
			typename: typename,
			typ:      types.Unalias(typename.Type()),
			methods:  mapx.NewXmap[string, *Function](),
		}

		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(node ast.Node) bool {
				f, ok := node.(*ast.FuncDecl)
				if !ok || f.Recv == nil {
					return true
				}
				s := reflectx.MustAssertType[*types.Signature](pkg.TypesInfo.TypeOf(f.Name))
				recv := deref(types.Unalias(s.Recv().Type()))
				typ := elem.Type()
				if !types.Identical(recv, typ) {
					return true
				}
				elem.methods.Store(f.Name.Name, &Function{
					node:     node,
					ident:    f.Name,
					function: pkg.TypesInfo.ObjectOf(f.Name).(*types.Func),
				})
				return true
			})
		}

		typenames.Append(elem)
	}

	typenames.Init()
	return typenames
}
