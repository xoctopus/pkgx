package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"sort"
	"strings"

	"github.com/xoctopus/x/misc/must"
	gopkg "golang.org/x/tools/go/packages"
)

// SourceOfNode helps for debugging node source code
func SourceOfNode(p *gopkg.Package, node ast.Node, trim bool) string {
	must.BeTrue(node != nil)

	// must a valid node
	must.BeTrue(node.Pos() != token.NoPos && node.End() != token.NoPos)
	pos := p.Fset.Position(node.Pos())
	must.BeTrue(pos != token.Position{})
	end := p.Fset.Position(node.End())
	must.BeTrue(end != token.Position{})

	b := bytes.NewBuffer(nil)
	must.NoError(printer.Fprint(b, p.Fset, node))

	filename := pos.Filename
	if trim {
		filename = strings.TrimPrefix(filename, p.Dir)
		filename = strings.TrimPrefix(filename, "/")
	}

	return fmt.Sprintf("%s\n"+
		"pos: %s:%d:%d\n"+
		"end: %s:%d:%d\n",
		strings.TrimSpace(b.String()),
		filename, pos.Line, pos.Column, filename, end.Line, end.Column,
	)
}

// func IdentNode(ident *ast.Ident) ast.Node {
// 	if ident != nil && ident.Obj != nil {
// 		if node, ok := ident.Obj.Decl.(ast.Node); ok {
// 			return node
// 		}
// 	}
// 	return nil
// }

func NodeOf(node ast.Node) Node {
	return Node{node.Pos(), node.End()}
}

// Node describes only pos and end as an ast.Node
type Node [2]token.Pos

func (p Node) Pos() token.Pos {
	return p[0]
}

func (p Node) End() token.Pos {
	return p[1]
}

func NodesOf[N ast.Node](e ...N) Nodes[N] {
	nodes := make(Nodes[N], len(e))
	copy(nodes, e)
	return nodes
}

type Nodes[N ast.Node] []N

func (ns Nodes[N]) Sort() Nodes[N] {
	sort.Slice(ns, func(i, j int) bool {
		return ns[i].Pos() < ns[j].Pos()
	})
	return ns
}

func Deref(t types.Type) types.Type {
	t = types.Unalias(t)
	for {
		if _t, ok := t.(*types.Pointer); ok {
			t = _t.Elem()
			continue
		}
		break
	}
	return t
}
