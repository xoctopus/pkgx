package pkgx

import (
	"go/types"

	"github.com/xoctopus/x/docx/v2"
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
	docs    map[string]*docx.Meta
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

func (t *TypeName) SetFieldDocs(docs map[string]*docx.Meta) {
	t.docs = docs
}

func (t *TypeName) GetFieldDocByName(name string) *docx.Meta {
	return t.docs[name]
}
