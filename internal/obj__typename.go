package internal

import (
	"go/types"

	"github.com/xoctopus/x/syncx"
)

func NewTypeName(obj Object[*types.TypeName]) *TypeName {
	return &TypeName{
		Object:  obj,
		methods: syncx.NewXmap[string, *Function](),
	}
}

type TypeName struct {
	Object[*types.TypeName]
	methods syncx.Map[string, *Function]
}

func (t *TypeName) Methods() syncx.Map[string, *Function] {
	return t.methods
}

func (t *TypeName) AddMethods(fns ...*Function) {
	for _, f := range fns {
		t.methods.Store(f.Name(), f)
	}
}

func (t *TypeName) Method(name string) *Function {
	f, _ := t.methods.Load(name)
	return f
}
