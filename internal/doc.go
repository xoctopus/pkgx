package internal

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

func ParseDocument(docs ...*ast.CommentGroup) *Doc {
	text := make([]string, 0, len(docs))
	for _, c := range docs {
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
			k = strings.TrimSpace(line[:idx])
			v = strings.TrimSpace(line[idx+1:])
		}
		d.tags[k] = append(d.tags[k], v)
	}
	if len(d.desc) == 0 && len(d.tags) == 0 {
		return d
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

	d.keys = maps.Keys(d.tags)
	sort.Strings(d.keys)
	return d
}

type Doc struct {
	tags map[string][]string
	keys []string
	desc []string
}

func (d *Doc) Tags() []string { return d.keys }

func (d *Doc) Values(tag string) []string { return d.tags[tag] }

func (d *Doc) Desc() string { return strings.Join(d.desc, " ") }

func (d *Doc) WithDesc(comments ...*ast.CommentGroup) *Doc {
	for _, c := range comments {
		if c != nil {
			for _, line := range strings.Split(c.Text(), "\n") {
				if text := strings.TrimSpace(line); len(text) > 0 {
					d.desc = append(d.desc, text)
				}
			}
		}
	}
	return d
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
