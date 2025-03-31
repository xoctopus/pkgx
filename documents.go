package pkgx

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"

	"golang.org/x/exp/maps"

	gopkg "golang.org/x/tools/go/packages"
)

func ParseAnnotation(text []string) *Annotation {
	anno := &Annotation{tags: make(map[string][]string)}
	for _, line := range text {
		if line[0] != '+' {
			anno.desc = append(anno.desc, line)
			continue
		}
		line = line[1:]
		k, v := "", ""
		if idx := strings.Index(line, "="); idx == -1 {
			k = line
		} else {
			k, v = strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:])
		}
		anno.tags[k] = append(anno.tags[k], v)
	}
	if len(anno.desc) == 0 && len(anno.tags) == 0 {
		return nil
	}
	anno.keys = maps.Keys(anno.tags)
	sort.Slice(anno.keys, func(i, j int) bool {
		return anno.keys[i] < anno.keys[j]
	})
	return anno
}

type Annotation struct {
	tags map[string][]string
	keys []string
	desc []string
}

func (c *Annotation) Tags() []string {
	return c.keys
}

func (c *Annotation) Values(tag string) []string {
	return c.tags[tag]
}

func (c *Annotation) Desc() []string {
	return c.desc
}

func (c *Annotation) String() string {
	b := strings.Builder{}
	b.WriteString("tags:[")
	for i, key := range c.Tags() {
		if i > 0 {
			b.WriteString("][")
		}
		b.WriteString(key)
		b.WriteString(":")
		b.WriteString(strings.Join(c.Values(key), ","))
	}
	b.WriteString("] desc:[")
	for i, line := range c.Desc() {
		if i > 0 {
			b.WriteString("][")
		}
		b.WriteString(line)
	}
	b.WriteString("]")
	return b.String()
}

func ParseDocuments(node ast.Node) []*Document {
	parse := func(node ast.Node, ident *ast.Ident, comments ...*ast.CommentGroup) *Document {
		text := make([]string, 0, len(comments))
		for _, c := range comments {
			if c == nil {
				continue
			}
			for _, line := range strings.Split(c.Text(), "\n") {
				if line = strings.TrimSpace(line); len(line) > 0 {
					text = append(text, line)
				}
			}

		}
		return &Document{
			node:  node,
			ident: ident,
			anno:  ParseAnnotation(text),
		}
	}

	var docs []*Document
	switch x := node.(type) {
	case *ast.GenDecl:
		for _, s := range x.Specs {
			switch sx := s.(type) {
			case *ast.ValueSpec:
				for _, ident := range sx.Names {
					docs = append(docs, parse(sx, ident, x.Doc, sx.Doc, sx.Comment))
				}
			case *ast.TypeSpec:
				docs = append(docs, parse(sx, sx.Name, x.Doc, sx.Doc, sx.Comment))
			}
		}
	case *ast.FuncDecl:
		docs = append(docs, parse(x, x.Name, x.Doc))
	}
	return docs
}

type Document struct {
	node  ast.Node
	ident *ast.Ident
	anno  *Annotation
}

func (d *Document) IsZero() bool {
	return d == nil || d.anno == nil
}

func (d *Document) Name() string {
	return d.ident.Name
}

func (d *Document) Node() ast.Node {
	return d.node
}

func (d *Document) Annotation() *Annotation {
	return d.anno
}

func (d *Document) String() string {
	return fmt.Sprintf("%s %s", d.Name(), d.Annotation())
}

func ScanDocuments(pkg *gopkg.Package) *Set[*Document] {
	d := NewSet[*Document](pkg)

	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			d.Append(ParseDocuments(decl)...)
		}
	}

	d.Init()
	return d
}
