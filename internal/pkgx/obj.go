package pkgx

import (
	"go/ast"
	"go/token"
	"go/types"
	"iter"
	"sort"

	"github.com/xoctopus/x/syncx"
)

// Exposer presents package level exposer excludes *types.Var
type Exposer interface {
	*types.Func | *types.Const | *types.TypeName
}

// Object defines parsed universal objects
type Object[U Exposer] interface {
	IsNil() bool
	Exposer() U
	Name() string
	Node() ast.Node
	Ident() *ast.Ident
	Doc() *Doc
	Type() types.Type
	TypeName() string
}

func NewObject[U Exposer](n ast.Node, i *ast.Ident, obj U, c *Doc) Object[U] {
	return &object[U]{node: n, id: i, u: obj, doc: c}
}

type object[U Exposer] struct {
	u    U
	node ast.Node
	id   *ast.Ident
	doc  *Doc
}

func (o *object[U]) IsNil() bool {
	return o == nil || o.node == nil
}

func (o *object[U]) Node() ast.Node {
	return o.node
}

func (o *object[U]) Name() string {
	if o.id != nil {
		return o.id.Name
	}
	return ""
}

func (o *object[U]) Ident() *ast.Ident {
	return o.id
}

func (o *object[U]) Exposer() U {
	return o.u
}

func (o *object[U]) Doc() *Doc {
	return o.doc
}

func (o *object[U]) Type() types.Type {
	if o.u == *new(U) {
		return nil
	}
	return any(o.u).(types.Object).Type()
}

func (o *object[U]) TypeName() string {
	switch t := any(o.u).(type) {
	case *types.TypeName:
		return t.Name()
	case *types.Const:
		if named, ok := t.Type().(*types.Named); ok {
			return named.Obj().Name()
		}
	}
	return ""
}

// Objects defines an interface for lookup and traverse Object by ast.Node or
// object name
type Objects[U Exposer, V Object[U]] interface {
	Len() int
	Nodes() iter.Seq[ast.Node]
	ExposerOf(ast.Node) U
	Exposers() iter.Seq[U]
	ElementOf(ast.Node) V
	Elements() iter.Seq[V]
	ElementByName(string) V
}

type ObjectsManager[U Exposer, V Object[U]] interface {
	Add(...V)
	Init(*token.FileSet)
	RangeNodes(func(Node, V) bool)
}

func NewObjects[U Exposer, V Object[U]]() Objects[U, V] {
	return &objects[U, V]{set: syncx.NewXmap[Node, V]()}
}

type objects[U Exposer, V Object[U]] struct {
	set   syncx.Map[Node, V]
	nodes []ast.Node
	vals  []V
}

func (s *objects[U, V]) Init(fileset *token.FileSet) {
	nodes := make(Nodes[ast.Node], 0)
	for node := range s.set.Range {
		nodes = append(nodes, node)
	}

	sort.Slice(nodes, func(i, j int) bool {
		pi, pj := fileset.Position(nodes[i].Pos()), fileset.Position(nodes[j].Pos())
		if pi.Filename == pj.Filename {
			return pi.Offset < pj.Offset
		}
		return pi.Filename < pj.Filename
	})

	for _, node := range nodes {
		e, _ := s.set.Load(NodeOf(node))
		s.nodes = append(s.nodes, e.Node())
		s.vals = append(s.vals, e)
	}
}

func (s *objects[U, V]) Len() int {
	return len(s.nodes)
}

func (s *objects[U, V]) Nodes() iter.Seq[ast.Node] {
	return func(yield func(ast.Node) bool) {
		for _, node := range s.nodes {
			if !yield(node) {
				return
			}
		}
	}
}

func (s *objects[U, V]) Add(elems ...V) {
	for _, e := range elems {
		if e.IsNil() {
			continue
		}
		s.set.LoadOrStore(NodeOf(e.Node()), e)
	}
}

func (s *objects[U, V]) ExposerOf(node ast.Node) U {
	if u, ok := s.set.Load(NodeOf(node)); ok {
		return u.Exposer()
	}
	return *new(U)
}

func (s *objects[U, V]) Exposers() iter.Seq[U] {
	return func(yield func(U) bool) {
		for _, v := range s.vals {
			if !yield(v.Exposer()) {
				return
			}
		}
	}
}

func (s *objects[U, V]) ElementOf(node ast.Node) V {
	e, _ := s.set.Load(NodeOf(node))
	return e
}

func (s *objects[U, V]) Elements() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range s.vals {
			if !yield(v) {
				return
			}
		}
	}
}

func (s *objects[U, V]) ElementByName(name string) (e V) {
	s.set.Range(func(_ Node, v V) bool {
		if v.Name() == name {
			e = v
			return false
		}
		return true
	})
	return
}

func (s *objects[U, V]) RangeNodes(f func(Node, V) bool) {
	s.set.Range(f)
}
