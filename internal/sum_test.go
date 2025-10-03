package pkgx_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/pkgx"
)

func TestLoadSumFile(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	cwd := filepath.Dir(filename)
	_ = os.RemoveAll(filepath.Join(cwd, "testdata", pkgx.SumFilename))

	defer func() {
		_ = os.RemoveAll(filepath.Join(cwd, "testdata", pkgx.SumFilename))
	}()

	t.Run("NilModule", func(t *testing.T) {
		NewWithT(t).Expect(pkgx.LoadSumFile(nil)).To(BeNil())
		p := u.Package("io")
		NewWithT(t).Expect(pkgx.LoadSumFile(p.Module())).To(BeNil())
	})

	mod := "github.com/xoctopus/pkgx/testdata"
	u := pkgx.NewPackages(mod)
	p := u.Package(mod)
	NewWithT(t).Expect(p).NotTo(BeNil())
	NewWithT(t).Expect(p.Module().Path).To(Equal(mod))
	t.Run("NoSumFile", func(t *testing.T) {
		NewWithT(t).Expect(pkgx.LoadSumFile(p.Module())).To(BeNil())
	})

	t.Run("Loaded", func(t *testing.T) {
		// xgo is not support go1.24
		// t.Run("FailedToOpenFile", func(t *testing.T) {
		// 	mock.Patch(os.OpenFile, func(string, int, os.FileMode) (*os.File, error) {
		// 		return nil, errors.New(t.Name())
		// 	})
		// 	NewWithT(t).Expect(u.Sum().Save().Error()).To(Equal(t.Name()))
		// })

		NewWithT(t).Expect(u.ModuleSum(mod).Save()).To(BeNil())

		s := pkgx.LoadSumFile(p.Module())
		NewWithT(t).Expect(s).NotTo(BeNil())
		NewWithT(t).Expect(s.Dir()).To(Equal(p.Module().Dir))

		h := s.Hash("github.com/xoctopus/pkgx/testdata/sub")
		NewWithT(t).Expect(len(h) > 0).To(BeTrue())
		h = s.Hash("github.com/xoctopus/pkgx/testdata")
		NewWithT(t).Expect(len(h) > 0).To(BeTrue())
	})

}
