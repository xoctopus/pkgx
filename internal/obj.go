package internal

import (
	"go/ast"
	"go/types"
	"iter"

	"github.com/xoctopus/x/mapx"
)

type Underlying interface {
	// TODO maybe has *types.Signature
	*types.Func | *types.Const | *types.TypeName
}

// Object defines parsed universal objects
type Object[U Underlying] interface {
	IsNil() bool
	Underlying() U
	Name() string
	Node() ast.Node
	Ident() *ast.Ident
	Doc() *Doc
	Type() types.Type
	TypeName() string
}

func NewObject[U Underlying](n ast.Node, i *ast.Ident, obj U, c *Doc) Object[U] {
	// TODO ensure ident
	return &object[U]{node: n, id: i, u: obj, doc: c}
}

type object[U Underlying] struct {
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

func (o *object[U]) Underlying() U {
	return o.u
}

func (o *object[U]) Doc() *Doc {
	return o.doc
}

func (o *object[U]) Type() types.Type {
	if o.u == *new(U) {
		return nil
	}

	switch x := any(o.u).(type) {
	// TODO if treat types.Signature as an Underlying
	// case types.Type:
	// 	return x
	default:
		return x.(types.Object).Type()
	}
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

// TODO type Signature struct{ Object[*types.Signature] }

// Objects defines an interface for lookup and traverse Object by ast.Node or
// object name
type Objects[U Underlying, V Object[U]] interface {
	Init()
	Len() int
	Nodes() iter.Seq[ast.Node]
	Add(...V)
	Underlying(ast.Node) U
	Underlyings() iter.Seq[U]
	Element(ast.Node) V
	Elements() iter.Seq[V]
	ElementByName(string) V
}

func NewObjects[U Underlying, V Object[U]]() Objects[U, V] {
	return &objects[U, V]{set: mapx.NewXmap[Node, V]()}
}

type objects[U Underlying, V Object[U]] struct {
	set   mapx.Map[Node, V]
	nodes []ast.Node
	vals  []V
}

func (s *objects[U, V]) Init() {
	size := mapx.Len(s.set)
	s.nodes = make([]ast.Node, size)
	s.vals = make([]V, size)

	nodes := make(Nodes[ast.Node], 0, size)
	for _, node := range mapx.Keys(s.set) {
		nodes = append(nodes, node)
	}
	nodes.Sort()

	for i, node := range nodes {
		e, _ := s.set.Load(NodeOf(node))
		s.nodes[i] = e.Node()
		s.vals[i] = e
	}
}

func (s *objects[U, V]) Len() int {
	return mapx.Len(s.set)
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

func (s *objects[U, V]) Underlying(node ast.Node) U {
	if u, ok := s.set.Load(NodeOf(node)); ok {
		return u.Underlying()
	}
	return *new(U)
}

func (s *objects[U, V]) Underlyings() iter.Seq[U] {
	return func(yield func(U) bool) {
		for _, v := range s.vals {
			if !yield(v.Underlying()) {
				return
			}
		}
	}
}

func (s *objects[U, V]) Element(node ast.Node) V {
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
