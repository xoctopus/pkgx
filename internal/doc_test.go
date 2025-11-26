package internal_test

import (
	"go/ast"
	"path/filepath"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/pkgx/internal"
)

// doc is document of `IntConstType`
var doc *internal.Doc

func init() {
	for _, f := range pkg.GoPackage().Syntax {
		ast.Inspect(f, func(node ast.Node) bool {
			if decl, ok := node.(*ast.GenDecl); ok {
				for _, spec := range decl.Specs {
					if spec, ok := spec.(*ast.TypeSpec); ok && spec.Name.Name == "IntConstType" {
						doc = internal.ParseDocument(decl.Doc, spec.Doc, spec.Comment)
						return false
					}
				}
			}
			return true
		})
	}
}

func TestParseDocument(t *testing.T) {
	t.Run("ParseDocumentForIntConstType", func(t *testing.T) {
		Expect(t, doc, NotBeNil[*internal.Doc]())
		Expect(t, doc.TagKeys(), Equal([]string{"key1", "key2", "key3", "key4"}))
		Expect(t, doc.Desc(), Equal([]string{
			"IntConstType defines a named constant type with integer underlying in a single `GenDecl`",
			"line1",
			"line2",
			"this is an inline comment",
		}))
		Expect(t, doc.TagValues("key3"), HaveLen[[]string](0))
		Expect(t, doc.TagValues("key4"), Equal([]string{"val_key4"}))
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
		Expect(t, doc.Tags(), Equal(map[string][]string{
			"key1": {"val_key1_1", "val_key1_2"},
			"key2": {"val_key2"},
			"key3": {},
			"key4": {"val_key4"},
		}))
	})

	t.Run("ParsePackageDocument", func(t *testing.T) {
		var doc *internal.Doc
		for _, f := range testdata.Syntax {
			if filepath.Base(testdata.Fset.File(f.Pos()).Name()) == "doc.go" {
				doc = internal.ParseDocument(f.Doc, f.Comments...)
			}
		}
		Expect(t, doc.String(), Equal("tags:[genx:apis][genx:enum][genx:model] desc:[Package testdata contains testdata for pkgx.][package desc following here][file comment here]"))
	})

	t.Run("NoDocument", func(t *testing.T) {
		Expect(t, internal.ParseDocument(nil).String(), Equal("tags: desc:"))
	})
}

func TestDocFields(t *testing.T) {
}
