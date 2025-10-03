package pkgx

import (
	"go/ast"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

func ParseDocument(comments ...*ast.CommentGroup) *Doc {
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

	d := &Doc{tags: make(map[string][]string)}
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
			k, v = strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:])
		}
		d.tags[k] = append(d.tags[k], v)
	}
	if len(d.desc) == 0 && len(d.tags) == 0 {
		return nil
	}
	d.keys = maps.Keys(d.tags)
	sort.Slice(d.keys, func(i, j int) bool {
		return d.keys[i] < d.keys[j]
	})
	return d
}

type Doc struct {
	tags map[string][]string
	keys []string
	desc []string
}

func (d *Doc) Tags() []string { return d.keys }

func (d *Doc) Values(tag string) []string { return d.tags[tag] }

func (d *Doc) Desc() []string { return d.desc }

func (d *Doc) Comment() string {
	if len(d.desc) > 0 {
		return strings.Join(d.desc, " ")
	}
	return ""
}

func (d *Doc) String() string {
	b := strings.Builder{}
	b.WriteString("tags:[")
	for i, key := range d.Tags() {
		if i > 0 {
			b.WriteString("][")
		}
		b.WriteString(key)
		b.WriteString(":")
		b.WriteString(strings.Join(d.Values(key), ","))
	}
	b.WriteString("] desc:[")
	for i, line := range d.Desc() {
		if i > 0 {
			b.WriteString("][")
		}
		b.WriteString(line)
	}
	b.WriteString("]")
	return b.String()
}
