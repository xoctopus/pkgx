package pkgx_test

import (
	"context"
	"go/types"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/pkgx"
)

func TestMustLoad(t *testing.T) {
	ctx := context.Background()
	ctx = pkgx.CtxLoadMode.With(ctx, pkgx.DefaultLoadMode)
	ExpectPanic[error](t, func() { pkgx.Load(ctx, "github.com/xoctopus/pkgx_test") })

	ctx = pkgx.CtxLoadTests.With(ctx, true)
	p := pkgx.Load(ctx, "github.com/xoctopus/pkgx/internal_test")
	Expect(t, p, NotBeNil[pkgx.Package]())
	Expect(t, p.ID(), NotEqual(p.Path()))
	Expect(t, p.WrapID(), HavePrefix("xwrap_"))

	ExpectPanic[error](t, func() { pkgx.Load(ctx, "example.com/a/b/c") })

	Expect(t, pkgx.Load(ctx, "io"), NotBeNil[pkgx.Package]())
}

func TestLookup(t *testing.T) {
	x, _ := pkgx.Lookup[*types.Named](context.Background(), "bytes", "Buffer")
	Expect(t, x, NotBeNil[*types.Named]())
	Expect(t, x.String(), Equal("bytes.Buffer"))
	x2 := pkgx.MustLookup[*types.Named](context.Background(), "bytes", "Buffer")
	Expect(t, x2, NotBeNil[*types.Named]())
	Expect(t, x2.String(), Equal(x.String()))

	t.Run("TypeNotMatch", func(t *testing.T) {
		x, _ := pkgx.Lookup[*types.Signature](context.Background(), "bytes", "Buffer")
		Expect(t, x, BeNil[*types.Signature]())
		ExpectPanic[error](t, func() {
			pkgx.MustLookup[*types.Signature](context.Background(), "bytes", "Buffer")
		})
	})

	t.Run("NoIdentifier", func(t *testing.T) {
		x, _ := pkgx.Lookup[types.Type](context.Background(), "bytes", "NoIdentifier")
		Expect(t, x, BeNil[types.Type]())
		ExpectPanic[error](t, func() {
			pkgx.MustLookup[types.Type](context.Background(), "bytes", "NoIdentifier")
		})
	})
}
