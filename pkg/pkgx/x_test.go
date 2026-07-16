package pkgx_test

import (
	"context"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/xoctopus/x/contextx"
	"github.com/xoctopus/x/docx/v2"
	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/pkgx/pkg/pkgx"
	_ "github.com/xoctopus/pkgx/testdata"
)

// init module dir
var (
	// make sure unit tests run without go workspace
	_ = os.Setenv("GOWORK", "off")

	testdata = "github.com/xoctopus/pkgx/testdata"
	version  = "v1.0.0" // in go.mod
	sub      = "github.com/xoctopus/pkgx/testdata/sub"
	cwd, _   = os.Getwd()
	dir      = filepath.Join(cwd, "..", "..", "testdata")

	u = NewPackages(
		contextx.Compose(
			CtxWorkdir.Carry(dir),
			CtxLoadTests.Carry(true),
		)(context.Background()),
		testdata,
	)
	pkg = u.Package(testdata)
)

func TestNewPackage(t *testing.T) {
	t.Run("Basics", func(t *testing.T) {
		Expect(t, pkg.Unwrap().Path(), Equal(testdata))
		Expect(t, pkg.Unwrap().Name(), Equal("testdata"))
		Expect(t, pkg.GoPackage().ID, Equal(testdata))
		Expect(t, pkg.GoPackage().Name, Equal("testdata"))
		Expect(t, pkg.Name(), Equal("testdata"))
		Expect(t, pkg.Path(), Equal(testdata))
		Expect(t, pkg.ID(), Equal(testdata))
		_ = pkg.Files()
		_ = pkg.FileSet()
		Expect(t, pkg.PackageDoc().Lines(), Equal([]string{
			"Package testdata contains testdata for pkgx.",
			"package desc following here",
		}))
		Expect(t, pkg.PackageDoc().Directives(), Equal([]string{
			"+genx:enum",
			"+genx:apis",
			"+genx:model",
		}))
	})

	t.Run("DifferentPathAndID", func(t *testing.T) {
		path := testdata + "_test"
		p := u.Package(path)
		Expect(t, p, NotBeNil[Package]())
		Expect(t, p.Path(), Equal(path))
		Expect(t, p.ID(), NotEqual(path))
	})

	t.Run("Module", func(t *testing.T) {
		Expect(t, pkg.GoModule().Path, Equal(testdata))
		Expect(t, pkg.PackageByPath(sub).GoModule().Path, Equal(testdata))
		Expect(t, pkg.PackageByPath("not/imported"), BeNil[Package]())

		Expect(t, pkg.PackageByPath("fmt").SourceDir(), Equal(""))
		Expect(t, pkg.SourceDir(), Equal(dir))
		Expect(t, pkg.SourceDir(), Equal(pkg.GoModule().Dir))
		Expect(t, pkg.PackageByPath(sub).SourceDir(), Equal(filepath.Join(dir, "sub")))
		Expect(t, u.ModuleSum(testdata).Dir(), Equal(pkg.GoModule().Dir))
	})

	t.Run("LookupAndEval", func(t *testing.T) {
		c := pkg.Constants().ElementByName("IntConstTypeValue1")
		n := pkg.TypeNames().ElementByName("TypeA")
		f := pkg.Functions().ElementByName("F")

		Expect(t, pkg.FieldDoc("Structure", "name").Lines(), Equal([]string{"name comments"}))
		Expect(t, pkg.FieldDoc("_", ""), BeNil[*docx.Meta]())
		Expect(t, pkg.Position(f.Node().Pos()).String(), Equal(filepath.Join(dir, "functions.go:22:1")))
		Expect(t, pkg.ObjectOf(f.Ident()).Name(), Equal("F"))

		_, err := pkg.Eval(nil)
		Expect(t, err, NotBeNil[error]())

		tav, _ := pkg.Eval(f.Ident())
		Expect(t, tav.Type.String(), Equal("func()"))
		tav, _ = pkg.Eval(n.Ident())
		Expect(t, tav.Type.String(), Equal("github.com/xoctopus/pkgx/testdata.TypeA"))
		tav, _ = pkg.Eval(c.Ident())
		Expect(t, tav.Type.String(), Equal("github.com/xoctopus/pkgx/testdata.IntConstType"))
		Expect(t, tav.Value.String(), Equal("1"))
	})
}

func ExamplePackage_Constants() {
	for e := range pkg.Constants().Elements() {
		fmt.Println(e.Name())
		fmt.Println(e.Doc().Lines())
		fmt.Printf("%s = %v\n", e.Name(), e.Value())
	}

	// Output:
	// IntConstTypeValue1
	// [IntConstTypeValue1 doc comment 1]
	// IntConstTypeValue1 = 1
	// IntConstTypeValue2
	// [IntConstTypeValue2 doc comment 2]
	// IntConstTypeValue2 = 2
	// IntConstTypeValue3
	// [comment 3]
	// IntConstTypeValue3 = 4
	// INT_STRING_ENUM__UNKNOWN
	// []
	// INT_STRING_ENUM__UNKNOWN = 0
	// INT_STRING_ENUM_A
	// [INT_STRING_ENUM_A has doc A]
	// INT_STRING_ENUM_A = 1
	// INT_STRING_ENUM_B
	// [has comment B]
	// INT_STRING_ENUM_B = 2
	// INT_STRING_ENUM_C
	// []
	// INT_STRING_ENUM_C = 3
	// Multi1
	// [multi ident will skip extract documents and nodes]
	// Multi1 = 1
}

func ExamplePackage_TypeNames() {
	for o := range pkg.TypeNames().Elements() {
		fmt.Println(o.Name())
		fmt.Println(o.Doc().Lines())

		fmt.Println(o.Ident().Name)
		methods := o.Methods().Keys()
		sort.Strings(methods)
		for _, name := range methods {
			m := o.Method(name)
			ref := ""
			if m.PtrRecv() {
				ref = "*"
			}
			fmt.Printf("%s%s.%s: %s\n", ref, o.Ident().Name, m.Name(), m.Type().String())
		}
		fmt.Println()
	}

	// Output:
	// IntConstType
	// [IntConstType defines a named constant type with integer underlying in a single `GenDecl` line1 line2 this is an inline comment]
	// IntConstType
	//
	// TypeA
	// [GenDecl defines 2 type, TypeA and TypeB TypeA doc line1 line2]
	// TypeA
	//
	// TypeB
	// [GenDecl defines 2 type, TypeA and TypeB TypeB doc line1 line2]
	// TypeB
	//
	// IntStringEnum
	// [IntStringEnum defines a named constant type with integer underlying as an enum type]
	// IntStringEnum
	//
	// Structure
	// [Structure is a struct type for testing line1 line2]
	// Structure
	// *Structure.Name: func() string
	// Structure.String: func() string
	// Structure.Value: func() any
	//
	// StructureAlias
	// [StructureAlias is an alias of Structure for testing]
	// StructureAlias
	//
	// Int
	// [type specs Int redefines int]
	// Int
	//
	// String
	// [type specs String redefines string]
	// String
	//
	// Float
	// [type specs Float alias of float64]
	// Float
	//
	// EachFieldHasComment
	// [EachFieldHasComment for field document]
	// EachFieldHasComment
}

func ExamplePackage_Functions() {
	for o := range pkg.Functions().Elements() {
		fmt.Println(o.Ident().Name)
		fmt.Println(o.Type())
		fmt.Println(o.Doc().Lines())
	}

	// Output:
	// Curry
	// func() func() int
	// [Curry function]
	// F
	// func()
	// [F a function list call expressions]
}

func ExamplePackages() {
	paths := make([]string, 0)
	fmt.Println("imported in company:")
	for path := range u.Packages {
		if strings.HasPrefix(path, "github.com/xoctopus") &&
			!(strings.HasSuffix(path, ".test") || strings.HasSuffix(path, "_test")) {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Println(path)
	}

	fmt.Println("directs")
	paths = paths[:0]
	for path := range u.Directs {
		if !(strings.HasSuffix(path, ".test") || strings.HasSuffix(path, "_test")) {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Println(path)
	}

	fmt.Println("modules")
	paths = paths[:0]
	for path := range u.Modules {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Println(path)
	}

	// Output:
	// imported in company:
	// github.com/xoctopus/pkgx/testdata
	// github.com/xoctopus/pkgx/testdata/sub
	// directs
	// github.com/xoctopus/pkgx/testdata
	// github.com/xoctopus/pkgx/testdata/sub
	// modules
	// github.com/xoctopus/pkgx/testdata
}

func TestWithWorkdir(t *testing.T) {
	t.Run("InModule", func(t *testing.T) {
		ctx := contextx.Compose(
			CtxWorkdir.Carry(filepath.Join(dir, "sub")),
			CtxLoadMode.Carry(DefaultLoadMode),
		)(context.Background())

		x := NewPackages(ctx, "github.com/xoctopus/pkgx/testdata/sub")
		p := x.Package("github.com/xoctopus/pkgx/testdata/sub")

		Expect(t, p.SourceDir(), Equal(filepath.Join(dir, "sub")))
	})

	t.Run("OutModule", func(t *testing.T) {
		dir := filepath.Join(cwd, "..", "..", "..", "internal", "devpkg", "consts")
		if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
			t.Skipf("skipping test because %s does not exist", dir)
		}

		ctx := CtxWorkdir.With(context.Background(), dir)
		pkgid := "github.com/xoctopus/internal/devpkg/consts"

		x := NewPackages(ctx, pkgid)
		p := x.Package(pkgid)
		Expect(t, p.SourceDir(), Equal(dir))
		o := p.Unwrap().Scope().Lookup("DoNotDelete")
		Expect(t, o, NotBeNil[types.Object]())
	})
}
