package pkgx

import (
	"go/ast"
	"iter"
	"sync"

	"github.com/xoctopus/x/misc/must"
	gopkg "golang.org/x/tools/go/packages"
)

var gNodes sync.Map

func NodeSeq(p *gopkg.Package) iter.Seq[ast.Node] {
	gNodes.LoadOrStore(p, sync.OnceValue(
		func() iter.Seq[ast.Node] {
			nodes := make([]ast.Node, 0)

			for _, f := range p.Syntax {
				ast.Inspect(f, func(n ast.Node) bool {
					if n != nil {
						nodes = append(nodes, n)
					}
					return true
				})
			}

			return func(yield func(ast.Node) bool) {
				for _, n := range nodes {
					if !yield(n) {
						return
					}
				}
			}
		},
	))

	seq, ok := gNodes.Load(p)
	must.BeTrue(ok)
	return seq.(func() iter.Seq[ast.Node])()
}
