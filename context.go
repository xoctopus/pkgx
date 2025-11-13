package pkgx

import (
	"context"
	"go/token"
	"sync"

	"github.com/xoctopus/x/contextx"
	gopkg "golang.org/x/tools/go/packages"
)

type PkgNamer interface {
	Package(id string) (name string)
}

const (
	LoadFiles       = gopkg.LoadFiles
	LoadImports     = gopkg.LoadImports
	LoadTypes       = gopkg.LoadTypes
	LoadSyntax      = gopkg.LoadSyntax
	LoadAllSyntax   = gopkg.LoadAllSyntax
	DefaultLoadMode = LoadAllSyntax | gopkg.NeedModule
)

var (
	ctxWorkdir   = contextx.NewT[string]()
	ctxLoadMode  = contextx.NewT[gopkg.LoadMode]()
	ctxPkgNamer  = contextx.NewT[PkgNamer]()
	ctxLogger    = contextx.NewT[func(string, ...any)]()
	ctxLoadTests = contextx.NewT[bool]()
	ctxFileset   = contextx.NewT[*token.FileSet]()
)

func WithLoadMode(ctx context.Context, mode gopkg.LoadMode) context.Context {
	return sync.OnceValue(func() context.Context {
		return ctxLoadMode.With(ctx, mode)
	})()
}

func LoadMode(ctx context.Context) gopkg.LoadMode {
	if mode, ok := ctxLoadMode.From(ctx); ok {
		return mode
	}
	return DefaultLoadMode
}

func WithWorkdir(ctx context.Context, workdir string) context.Context {
	return sync.OnceValue(func() context.Context {
		return ctxWorkdir.With(ctx, workdir)
	})()
}

func Workdir(ctx context.Context) string {
	if workdir, ok := ctxWorkdir.From(ctx); ok {
		return workdir
	}
	return ""
}

func WithNamer(ctx context.Context, namer PkgNamer) context.Context {
	return sync.OnceValue(func() context.Context {
		return ctxPkgNamer.With(ctx, namer)
	})()
}

func PackageName(ctx context.Context, p Package) string {
	if namer, ok := ctxPkgNamer.From(ctx); ok && namer != nil {
		return namer.Package(p.ID())
	}
	return p.Name()
}

func WithLogger(ctx context.Context, l func(string, ...any)) context.Context {
	return sync.OnceValue(func() context.Context {
		return ctxLogger.With(ctx, l)
	})()
}

func Logger(ctx context.Context) func(string, ...any) {
	if l, ok := ctxLogger.From(ctx); ok {
		return l
	}
	return nil
}

func WithTests(ctx context.Context) context.Context {
	return sync.OnceValue(func() context.Context {
		return ctxLoadTests.With(ctx, true)
	})()
}

func LoadTests(ctx context.Context) bool {
	_, ok := ctxLoadTests.From(ctx)
	return ok
}

func WithFileset(ctx context.Context, fileset *token.FileSet) context.Context {
	return sync.OnceValue(func() context.Context {
		return ctxFileset.With(ctx, fileset)
	})()
}

func Fileset(ctx context.Context) *token.FileSet {
	v, _ := ctxFileset.From(ctx)
	return v
}

func Config(ctx context.Context) *gopkg.Config {
	return &gopkg.Config{
		Fset:  Fileset(ctx),
		Mode:  LoadMode(ctx),
		Logf:  Logger(ctx),
		Dir:   Workdir(ctx),
		Tests: LoadTests(ctx),
	}
}
