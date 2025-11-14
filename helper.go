package pkgx

import (
	"context"
	"errors"
	"fmt"
	"go/types"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/syncx"
	gopkg "golang.org/x/tools/go/packages"

	"github.com/xoctopus/pkgx/internal"
)

func Lookup[T types.Type](ctx context.Context, path, x string) (T, bool) {
	p := Load(ctx, path)
	o := p.Unwrap().Scope().Lookup(x)
	if o != nil {
		t, ok := o.Type().(T)
		return t, ok
	}
	return *new(T), false
}

func MustLookup[T types.Type](ctx context.Context, path, name string) T {
	if t, ok := Lookup[T](ctx, path, name); ok {
		return t
	}
	panic(fmt.Errorf("expect lookup `%s.%s` as %T", path, name, reflect.TypeFor[T]()))
}

var gPkgs = syncx.NewXmap[string, func() Package]()

func load(ctx context.Context, path string) *gopkg.Package {
	_path := path
	if strings.HasSuffix(path, "_test") {
		path = strings.TrimSuffix(_path, "_test")
	}

	pkgs, err := gopkg.Load(Config(ctx), path)
	must.NoErrorF(err, "failed to load %s", path)
	must.BeTrueF(len(pkgs) > 0, "no packages loaded")
	must.NoErrorF(errors.Join(
		slices.Collect(func(yield func(error) bool) {
			for _, x := range pkgs[0].Errors {
				yield(x)
			}
		})...,
	), "failed to load %s", path)

	var pkg *gopkg.Package
	for i := range pkgs {
		if pkgs[i].PkgPath == _path {
			pkg = pkgs[i]
			break
		}
	}
	must.BeTrueF(pkg != nil, "package %s not loaded", _path)
	return pkg
}

func AsPackage(pkg *gopkg.Package) Package {
	return newx(pkg)
}

func Load(ctx context.Context, path string) Package {
	path = internal.NewWrapper().Unwrap(path)
	f, _ := gPkgs.LoadOrStore(path, sync.OnceValue(func() Package {
		return AsPackage(load(ctx, path))
	}))
	return f()
}
