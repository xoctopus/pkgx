package internal_test

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	. "github.com/onsi/gomega"
	. "github.com/xoctopus/pkgx/internal"
	"github.com/xoctopus/x/misc/must"
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
		code := SourceOfNode(testdata, node, cwd)
		fmt.Println(code)
	}

	// Output:
	// IntConstType int	// this is an inline comment
	// pos: ../testdata/documents.go:16:6
	// end: ../testdata/documents.go:16:22
	//
	// // IntConstTypeValue1 doc
	// IntConstTypeValue1 IntConstType = iota + 1	// comment 1
	// pos: ../testdata/documents.go:20:2
	// end: ../testdata/documents.go:20:44
	//
	// // TypeA doc
	// // line1
	// // line2
	// // +tag1=val1_1
	// // +tag1=val1_2
	// TypeA int
	// pos: ../testdata/documents.go:34:2
	// end: ../testdata/documents.go:34:11
	//
	// // TypeB doc
	// // line1
	// // line2
	// // +tag1=val1_1
	// // +tag1=val1_2
	// TypeB string
	// pos: ../testdata/documents.go:40:2
	// end: ../testdata/documents.go:40:14
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

	NewWithT(t).Expect(typ).NotTo(BeNil())
	ptr := types.NewPointer(typ)
	NewWithT(t).Expect(types.Identical(typ, Deref(ptr))).To(BeTrue())
}
