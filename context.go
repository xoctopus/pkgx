package pkgx

import (
	"context"

	"github.com/xoctopus/x/contextx"
)

var ctxWorkdir = contextx.NewT[string]()

func WithWorkdir(ctx context.Context, workdir string) context.Context {
	if _, ok := ctxWorkdir.From(ctx); ok {
		return ctx
	}
	return ctxWorkdir.With(ctx, workdir)
}
