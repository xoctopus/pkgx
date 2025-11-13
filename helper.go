package pkgx

import (
	"context"
	"go/types"
	"reflect"

	"github.com/pkg/errors"
)

func Lookup[T types.Type](ctx context.Context, path, name string) (T, bool) {
	p := NewPackages(ctx).Package(path)
	if p != nil {
		o := p.Unwrap().Scope().Lookup(name)
		if o != nil {
			t, ok := o.Type().(T)
			return t, ok
		}
	}
	return *new(T), false
}

func MustLookup[T types.Type](ctx context.Context, path, name string) T {
	if t, ok := Lookup[T](ctx, path, name); ok {
		return t
	}
	panic(errors.Errorf("expect lookup `%s.%s` as %T", path, name, reflect.TypeFor[T]()))
}
