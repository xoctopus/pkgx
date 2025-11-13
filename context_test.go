package pkgx_test

import (
	"context"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/pkgx"
)

type namer struct{}

func (*namer) Package(p string) string {
	if p == "io" {
		return "std_io"
	}
	return p
}

func TestConfig(t *testing.T) {
	ctx := context.Background()

	defv := Config(ctx)
	Expect(t, defv.Tests, BeFalse())
	Expect(t, defv.Mode, Equal(DefaultLoadMode))
	Expect(t, defv.Dir, Equal(""))
	Expect(t, defv.Logf, BeNil[func(string, ...any)]())
	Expect(t, defv.Fset, BeNil[*token.FileSet]())

	Expect(t, PackageName(ctx, Load(ctx, "io")), Equal("io"))
	Expect(t, PackageName(ctx, Load(ctx, "bytes")), Equal("bytes"))

	workdir := filepath.Join(os.Getenv("GOROOT"), "src")

	ctx = WithWorkdir(ctx, workdir)
	ctx = WithTests(ctx)
	ctx = WithLoadMode(ctx, LoadImports)
	ctx = WithLogger(ctx, func(string, ...any) {})
	ctx = WithNamer(ctx, &namer{})
	ctx = WithFileset(ctx, token.NewFileSet())
	defv = Config(ctx)
	Expect(t, defv.Tests, BeTrue())
	Expect(t, defv.Mode, Equal(LoadImports))
	Expect(t, defv.Dir, Equal(workdir))
	Expect(t, defv.Logf, NotBeNil[func(string, ...any)]())
	Expect(t, defv.Fset, NotBeNil[*token.FileSet]())

	Expect(t, PackageName(ctx, Load(ctx, "io")), Equal("std_io"))
	Expect(t, PackageName(ctx, Load(ctx, "bytes")), Equal("bytes"))
}
