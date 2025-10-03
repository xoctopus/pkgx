package internal_test

import (
	"go/ast"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/pkgx/internal"
)

func TestParseDocument(t *testing.T) {
	t.Run("ParseDocumentForIntConstType", func(t *testing.T) {
		var doc *internal.Doc
		for _, f := range testdata.Syntax {
			ast.Inspect(f, func(node ast.Node) bool {
				if decl, ok := node.(*ast.GenDecl); ok {
					for _, spec := range decl.Specs {
						if spec, ok := spec.(*ast.TypeSpec); ok && spec.Name.Name == "IntConstType" {
							doc = internal.ParseDocument(decl.Doc, spec.Doc).WithDesc(spec.Comment)
							return false
						}
					}
				}
				return true
			})
		}

		Expect(t, doc, NotBeNil[*internal.Doc]())
		Expect(t, doc.Tags(), Equal([]string{"key1", "key2", "key3", "key4"}))
		Expect(t, doc.Desc(), Equal(
			"IntConstType defines a named constant type with integer underlying in a single `GenDecl` "+
				"line1 "+
				"line2 "+
				"this is an inline comment",
		))
		Expect(t, doc.Values("key3"), HaveLen[[]string](0))
		Expect(t, doc.Values("key4"), Equal([]string{"val_key4"}))
		Expect(t, doc.String(), Equal(
			"tags:"+
				"[key1:val_key1_1,val_key1_2]"+
				"[key2:val_key2]"+
				"[key3]"+
				"[key4:val_key4] "+
				"desc:"+
				"[IntConstType defines a named constant type with integer underlying in a single `GenDecl`]"+
				"[line1]"+
				"[line2]"+
				"[this is an inline comment]"))
	})

	t.Run("NoDocument", func(t *testing.T) {
		doc := internal.ParseDocument()
		Expect(t, doc.String(), Equal("tags: desc:"))
	})
}
