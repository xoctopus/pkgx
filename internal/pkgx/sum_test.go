package pkgx_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/xoctopus/x/misc/must"
	. "github.com/xoctopus/x/testx"
	gopkg "golang.org/x/tools/go/packages"

	"github.com/xoctopus/pkgx/internal/pkgx"
)

var (
	testdata *gopkg.Package // = "github.com/xoctopus/pkgx/testdata"
	sub      *gopkg.Package
	cwd      string
)

func init() {
	_ = os.Setenv("GOWORK", "off")

	_, filename, _, _ := runtime.Caller(0)
	cwd = filepath.Dir(filename)

	pkgs, err := gopkg.Load(&gopkg.Config{
		Dir:  filepath.Join(cwd, "..", "..", "testdata"),
		Mode: gopkg.LoadMode(0b11111111111111111),
	}, "github.com/xoctopus/pkgx/testdata")
	must.NoError(err)
	must.BeTrue(len(pkgs) == 1)
	testdata = pkgs[0]

	pkgs, err = gopkg.Load(&gopkg.Config{
		Mode: gopkg.LoadMode(0b11111111111111111),
	}, "github.com/xoctopus/pkgx/testdata/sub")
	must.NoError(err)
	must.BeTrue(len(pkgs) == 1)
	sub = pkgs[0]
}

func TestLoadSumFile(t *testing.T) {
	filename := filepath.Join(cwd, "../../testdata", pkgx.SumFilename)
	_ = os.RemoveAll(filename)

	t.Run("NoModule", func(t *testing.T) {
		Expect(t, pkgx.LoadSumFile(nil), BeNil[pkgx.Sum]())

		pkgs, err := gopkg.Load(nil, "io")
		Expect(t, err, BeNil[error]())
		Expect(t, pkgs, HaveLen[[]*gopkg.Package](1))
		Expect(t, pkgx.LoadSumFile(pkgs[0].Module), BeNil[pkgx.Sum]())
	})

	path := "github.com/xoctopus/pkgx/testdata"
	Expect(t, testdata.Module.Path, Equal(path))

	t.Run("NoSumFile", func(t *testing.T) {
		Expect(t, pkgx.LoadSumFile(testdata.Module), BeNil[pkgx.Sum]())
	})

	sum := pkgx.NewSum(testdata.Module.Dir)
	Expect(t, sum.Dir(), Equal(filepath.Dir(filename)))

	t.Run("AddPackagesHashes", func(t *testing.T) {
		sum.Add(testdata)
		h := sum.Hash(testdata.ID)
		Expect(t, h, NotEqual(""))

		sum.Add(sub)
		h = sum.Hash(sub.ID)
		Expect(t, h, NotEqual(""))
	})

	t.Run("SaveAndLoad", func(t *testing.T) {
		Expect(t, sum.Save(), Succeed())

		sum2 := pkgx.LoadSumFile(testdata.Module)
		Expect(t, sum2, NotBeNil[pkgx.Sum]())
		Expect(t, sum2.Hash(testdata.ID), Equal(sum.Hash(testdata.ID)))
		Expect(t, sum2.Hash(sub.ID), Equal(sum.Hash(sub.ID)))

		t.Run("FailedLoad", func(t *testing.T) {
			t.Run("FailedToOpenFile", func(t *testing.T) {
				info, _ := os.Stat(sum.Dir())
				mode := info.Mode()

				_ = os.Chmod(sum.Dir(), 0000)
				defer func() {
					_ = os.Chmod(sum.Dir(), mode)
				}()

				Expect(t, sum.Save(), ErrorContains("permission denied"))
			})
		})
	})
}
