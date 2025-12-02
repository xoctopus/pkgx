package pkgx

import (
	"context"
	"go/token"

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
	CtxWorkdir   = contextx.NewT[string](contextx.WithDefault(""))
	CtxLoadMode  = contextx.NewT[gopkg.LoadMode](contextx.WithDefault(DefaultLoadMode))
	CtxLogger    = contextx.NewT[func(string, ...any)](contextx.WithDefault[func(string, ...any)](nil))
	CtxLoadTests = contextx.NewT[bool](contextx.WithDefault(false))
	CtxFileset   = contextx.NewT[*token.FileSet](contextx.WithDefault[*token.FileSet](nil))
	CtxPkgNamer  = contextx.NewT[PkgNamer]()
)

func PackageName(ctx context.Context, p Package) string {
	if namer, ok := CtxPkgNamer.From(ctx); ok && namer != nil {
		return namer.Package(p.ID())
	}
	return p.Name()
}

func Config(ctx context.Context) *gopkg.Config {
	return &gopkg.Config{
		Fset:  CtxFileset.MustFrom(ctx),
		Mode:  CtxLoadMode.MustFrom(ctx),
		Logf:  CtxLogger.MustFrom(ctx),
		Dir:   CtxWorkdir.MustFrom(ctx),
		Tests: CtxLoadTests.MustFrom(ctx),
	}
}
