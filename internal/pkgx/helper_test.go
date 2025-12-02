package pkgx_test

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	"github.com/xoctopus/x/misc/must"
	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/pkgx/internal/pkgx"
)

func ExampleSourceOfNode() {
	var nodes []ast.Node

	for n := range NodeSeq(testdata) {
		if x, ok := n.(*ast.TypeSpec); ok {
			if x.Name != nil && (x.Name.Name == "TypeB" ||
				x.Name.Name == "TypeA" || x.Name.Name == "IntConstType") {
				nodes = append(nodes, n)
			}
		}
		if x, ok := n.(*ast.ValueSpec); ok {
			if len(x.Names) >= 1 && x.Names[0].Name == "IntConstTypeValue1" {
				nodes = append(nodes, n)
			}
		}
		if len(nodes) >= 4 {
			break
		}
	}

	for _, node := range NodesOf(nodes...).Sort() {
		n := NodeOf(node)
		must.BeTrue(n.Pos() == node.Pos() && n.End() == node.End())
		code := SourceOfNode(testdata, node, true)
		fmt.Println(code)
	}

	// Output:
	// IntConstType int	// this is an inline comment
	// pos: documents.go:18:6
	// end: documents.go:18:22
	//
	// // IntConstTypeValue1 doc
	// IntConstTypeValue1 IntConstType = iota + 1	// comment 1
	// pos: documents.go:22:2
	// end: documents.go:22:44
	//
	// // TypeA doc
	// // line1
	// // line2
	// // +tag1=val1_1
	// // +tag1=val1_2
	// TypeA int
	// pos: documents.go:36:2
	// end: documents.go:36:11
	//
	// // TypeB doc
	// // line1
	// // line2
	// // +tag1=val1_1
	// // +tag1=val1_2
	// TypeB string
	// pos: documents.go:42:2
	// end: documents.go:42:14
}

func TestDeref(t *testing.T) {
	var typ types.Type
	for _, f := range testdata.Syntax {
		ast.Inspect(f, func(node ast.Node) bool {
			if decl, ok := node.(*ast.GenDecl); ok {
				if decl.Tok == token.TYPE {
					for _, s := range decl.Specs {
						if sp := s.(*ast.TypeSpec); sp.Name.Name == "Structure" {
							typ = testdata.TypesInfo.Types[sp.Type].Type
							return false
						}
					}
				}
			}
			return true
		})
	}

	Expect(t, typ, NotBeNil[types.Type]())
	ptr := types.NewPointer(typ)
	Expect(t, types.Identical(typ, Deref(ptr)), BeTrue())
}
