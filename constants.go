package pkgx

import (
	"go/ast"
	"go/constant"
	"go/types"

	"github.com/xoctopus/x/misc/must"
	gopkg "golang.org/x/tools/go/packages"
)

type Constant struct {
	node  ast.Node
	ident *ast.Ident
	value *types.Const
}

func (c *Constant) IsZero() bool {
	return c == nil || c.value == nil
}

func (c *Constant) Name() string {
	return c.ident.Name
}

func (c *Constant) Node() ast.Node {
	return c.node
}

func (c *Constant) Value() constant.Value {
	return c.value.Val()
}

func ScanConstants(pkg *gopkg.Package) *Set[*Constant] {
	constants := NewSet[*Constant](pkg)

	for ident, def := range pkg.TypesInfo.Defs {
		value, ok := def.(*types.Const)
		if !ok {
			continue
		}
		node := IdentNode(ident)
		must.BeTrueWrap(node != nil, "%s: %s", ident.Name, SourceOfNode(pkg, node))
		constants.Append(&Constant{
			node:  node,
			ident: ident,
			value: value,
		})
	}

	constants.Init()
	return constants
}
