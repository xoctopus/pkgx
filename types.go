package pkgx

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"log/slog"
	"path/filepath"
	"sort"
	"sync"

	"github.com/xoctopus/x/mapx"
	"github.com/xoctopus/x/misc/must"
	gopkg "golang.org/x/tools/go/packages"
)

// Element is a types union for source element
type Element interface {
	IsZero() bool
	Name() string
	Node() ast.Node
}

type ElementSet[T Element] interface {
	Init()
	Values() []T
	Append(...T)
	ValueByNode(ast.Node) T
	ValueByName(string) T
	Range(func(T) bool)
}

func NewSet[T Element](pkg *gopkg.Package) *Set[T] {
	return &Set[T]{
		pkg:      pkg,
		elements: mapx.NewXmap[Node, T](),
		once:     &sync.Once{},
	}
}

type Set[T Element] struct {
	pkg      *gopkg.Package
	elements mapx.Map[Node, T]
	once     *sync.Once
	nodes    []ast.Node
	values   []T
}

func (s *Set[T]) Init() {
	s.once.Do(func() {
		nodes := mapx.Keys(s.elements)
		sort.Sort(Nodes(nodes))

		for _, node := range nodes {
			v, _ := s.elements.Load(node)
			s.values = append(s.values, v)
			s.nodes = append(s.nodes, v.Node())
		}
	})
}

func (s *Set[T]) Values() []T {
	return s.values
}

func (s *Set[T]) Append(elements ...T) {
	for _, v := range elements {
		if v.IsZero() {
			continue
		}
		prev, exists := s.elements.LoadOrStore(NodeOf(v.Node()), v)
		if exists {
			slog.Debug(
				"syntax node conflict",
				slog.String("prev", PositionOf(s.pkg, prev.Node()).String()),
				slog.String("curr", PositionOf(s.pkg, v.Node()).String()),
			)
		}
	}
}

func (s *Set[T]) ValueByNode(node ast.Node) (e T) {
	e, _ = s.elements.Load(NodeOf(node))
	return
}

func (s *Set[T]) ValueByName(name string) (e T) {
	s.elements.Range(func(k Node, v T) bool {
		if v.Name() == name {
			e = v
			return false
		}
		return true
	})
	return
}

func (s *Set[T]) Range(f func(v T) bool) {
	s.elements.Range(func(_ Node, v T) bool {
		return f(v)
	})
}

func PositionOf(pkg *gopkg.Package, node ast.Node) *Position {
	pos := pkg.Fset.Position(node.Pos())
	end := pkg.Fset.Position(node.End())
	return &Position{dir: pkg.Dir, pos: pos, end: end}
}

// Position describes a node position in package fileset
type Position struct {
	dir      string
	filename string
	pos      token.Position
	end      token.Position
}

func (p *Position) WithCwd(cwd string) *Position {
	p.dir, _ = filepath.Rel(cwd, p.dir)
	must.BeTrue(p.dir != "")
	return p
}

func (p *Position) String() string {
	filename := p.Filename()
	return fmt.Sprintf(
		"%s:%d:%d\n%s:%d:%d",
		filename, p.pos.Line, p.pos.Column, filename, p.end.Line, p.end.Column,
	)
}

func (p *Position) Filename() string {
	return filepath.Join(p.dir, filepath.Base(p.pos.Filename))
}

func IdentNode(ident *ast.Ident) ast.Node {
	if ident.Obj == nil {
		return nil
	}
	node, _ := ident.Obj.Decl.(ast.Node)
	return node
}

func NodeOf(node ast.Node) Node {
	return Node{node.Pos(), node.End()}
}

// Node describes pos and end of an ast.Node
type Node [2]token.Pos

func (p Node) Less(p2 Node) bool {
	if p.Pos() == p2.Pos() {
		return p.End() < p2.End()
	}
	return p.Pos() < p2.Pos()
}

func (p Node) Pos() token.Pos {
	return p[0]
}

func (p Node) End() token.Pos {
	return p[1]
}

type Nodes []Node

func (v Nodes) Len() int {
	return len(v)
}

func (v Nodes) Less(i, j int) bool {
	return v[i].Less(v[j])
}

func (v Nodes) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func deref(t types.Type) types.Type {
	for {
		if _t, ok := t.(*types.Pointer); ok {
			t = _t.Elem()
			continue
		}
		break
	}
	return t
}

func SourceOfNode(p *gopkg.Package, node ast.Node) string {
	b := bytes.NewBuffer(nil)
	must.NoError(printer.Fprint(b, p.Fset, node))
	return fmt.Sprintf("%s\n%s", b.String(), PositionOf(p, node))
}
