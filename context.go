package pkgx

import (
	"context"

	"github.com/xoctopus/x/contextx"
	gopkg "golang.org/x/tools/go/packages"
)

var (
	ctxWorkdir  = contextx.NewT[string]()
	ctxLoadMode = contextx.NewT[gopkg.LoadMode]()
)

func WithWorkdir(ctx context.Context, workdir string) context.Context {
	if _, ok := ctxWorkdir.From(ctx); ok {
		return ctx
	}
	return ctxWorkdir.With(ctx, workdir)
}

const (
	LoadFiles       = gopkg.LoadFiles
	LoadImports     = gopkg.LoadImports
	LoadTypes       = gopkg.LoadTypes
	LoadSyntax      = gopkg.LoadSyntax
	LoadAllSyntax   = gopkg.LoadAllSyntax
	DefaultLoadMode = LoadAllSyntax | gopkg.NeedModule
)

func WithLoadMode(ctx context.Context, mode gopkg.LoadMode) context.Context {
	if _, ok := ctxLoadMode.From(ctx); ok {
		return ctx
	}
	return ctxLoadMode.With(ctx, mode)
}
