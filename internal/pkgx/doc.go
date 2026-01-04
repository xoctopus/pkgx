package pkgx

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"

	"github.com/xoctopus/x/slicex"
)

func ParseDocument(doc *ast.CommentGroup, comments ...*ast.CommentGroup) *Doc {
	docs := make([]*ast.Comment, 0)
	for _, d := range append(comments, doc) {
		if d != nil {
			for _, c := range d.List {
				if c != nil {
					docs = append(docs, c)
				}
			}
		}
	}
	if len(docs) == 0 {
		return &Doc{tags: make(map[string][]string)}
	}

	docs = slicex.Unique(docs)
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Pos() < docs[j].Pos()
	})

	d := &Doc{
		tags:         make(map[string][]string),
		CommentGroup: &ast.CommentGroup{List: docs},
	}

	text := make([]string, 0, len(docs))
	for _, c := range docs {
		for _, line := range strings.Split(c.Text, "\n") {
			line = strings.TrimPrefix(line, "/*")
			line = strings.TrimPrefix(line, "//")
			line = strings.TrimSuffix(line, "*/")
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				text = append(text, line)
			}
		}
	}

	for _, line := range text {
		if line[0] != '+' {
			d.desc = append(d.desc, line)
			continue
		}
		line = line[1:]
		k, v := "", ""
		if idx := strings.Index(line, "="); idx == -1 {
			k = line
		} else {
			k = strings.TrimSpace(line[:idx])
			v = strings.TrimSpace(line[idx+1:])
		}
		d.tags[k] = append(d.tags[k], v)
	}
	if len(d.desc) == 0 && len(d.tags) == 0 {
		return &Doc{tags: make(map[string][]string)}
	}
	for tag, vals := range d.tags {
		trimmed := make([]string, 0, len(vals))
		for _, v := range vals {
			if v != "" {
				trimmed = append(trimmed, v)
			}
		}
		d.tags[tag] = trimmed
	}

	for key := range d.tags {
		d.keys = append(d.keys, key)
	}
	sort.Strings(d.keys)
	return d
}

var DefaultDoc = &Doc{tags: make(map[string][]string), CommentGroup: &ast.CommentGroup{}}

type Doc struct {
	tags map[string][]string
	keys []string // sorted
	desc []string
	*ast.CommentGroup
}

func (d *Doc) Tags() map[string][]string {
	return d.tags
}

func (d *Doc) TagKeys() []string {
	return d.keys
}

func (d *Doc) TagValues(tag string) []string {
	return d.tags[tag]
}

func (d *Doc) Desc() []string {
	return d.desc
}

func (d *Doc) String() string {
	tags := make([]string, 0, len(d.tags))
	for _, k := range d.keys {
		if vs := d.tags[k]; len(vs) > 0 {
			tags = append(tags, fmt.Sprintf("[%s:%s]", k, strings.Join(vs, ",")))
		} else {
			tags = append(tags, fmt.Sprintf("[%s]", k))
		}
	}

	desc := make([]string, 0, len(d.desc))
	for _, v := range d.desc {
		desc = append(desc, fmt.Sprintf("[%s]", v))
	}
	return fmt.Sprintf("tags:%s desc:%s", strings.Join(tags, ""), strings.Join(desc, ""))
}
