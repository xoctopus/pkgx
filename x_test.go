package pkgx_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/pkgx"
	_ "github.com/xoctopus/pkgx/testdata"
)

var (
	module   = "github.com/xoctopus/pkgx"
	testdata = "github.com/xoctopus/pkgx/testdata"
	sub      = "github.com/xoctopus/pkgx/testdata/sub"
	cwd, _   = os.Getwd()

	u   = NewPackages(context.Background(), module, testdata)
	pkg = u.Package(testdata)
)

func TestNewPackage(t *testing.T) {
	Expect(t, pkg.Unwrap().Path(), Equal(testdata))
	Expect(t, pkg.Unwrap().Name(), Equal("testdata"))
	Expect(t, pkg.GoPackage().PkgPath, Equal(testdata))
	Expect(t, pkg.GoPackage().Name, Equal("testdata"))

	Expect(t, pkg.GoModule().Path, Equal(testdata))
	Expect(t, pkg.PackageByPath(sub).GoModule().Path, Equal(testdata))

	Expect(t, pkg.PackageByPath("not/imported"), BeNil[Package]())

	Expect(t, pkg.PackageByPath("fmt").SourceDir(), Equal(""))
	Expect(t, pkg.SourceDir(), Equal(filepath.Join(cwd, "testdata")))
	Expect(t, pkg.SourceDir(), Equal(pkg.GoModule().Dir))
	Expect(t, pkg.PackageByPath(sub).SourceDir(), Equal(filepath.Join(cwd, "testdata", "sub")))

	c := pkg.Constants().ElementByName("IntConstTypeValue1")
	n := pkg.TypeNames().ElementByName("TypeA")
	f := pkg.Functions().ElementByName("F")

	Expect(t, pkg.Position(f.Node().Pos()).String(), HaveSuffix("testdata/functions.go:22:1"))
	Expect(t, pkg.ObjectOf(f.Ident()).Name(), Equal("F"))

	Expect(t, u.ModuleSum(testdata).Dir(), Equal(pkg.GoModule().Dir))

	_, err := pkg.Eval(nil)
	Expect(t, err, NotBeNil[error]())

	tav, _ := pkg.Eval(f.Ident())
	Expect(t, tav.Type.String(), Equal("func()"))
	tav, _ = pkg.Eval(n.Ident())
	Expect(t, tav.Type.String(), Equal("github.com/xoctopus/pkgx/testdata.TypeA"))
	tav, _ = pkg.Eval(c.Ident())
	Expect(t, tav.Type.String(), Equal("github.com/xoctopus/pkgx/testdata.IntConstType"))
	Expect(t, tav.Value.String(), Equal("1"))

	// TODO
	_ = pkg.Files()
	_ = pkg.FileSet()
	_ = pkg.Doc()
}

func ExamplePackage_Constants() {
	for e := range pkg.Constants().Elements() {
		fmt.Println(e.Doc())
		fmt.Printf("%s = %v\n", e.Name(), e.Value())
	}

	// Output:
	// tags: desc:[IntConstTypeValue1 doc][comment 1]
	// IntConstTypeValue1 = 1
	// tags: desc:[IntConstTypeValue2 doc][comment 2]
	// IntConstTypeValue2 = 2
	// tags: desc:[placeholder]
	// _ = 3
	// tags: desc:[comment 3]
	// IntConstTypeValue3 = 4
	// tags: desc:
	// INT_STRING_ENUM__UNKNOWN = 0
	// tags: desc:[INT_STRING_ENUM_A has doc A]
	// INT_STRING_ENUM_A = 1
	// tags: desc:[has comment B]
	// INT_STRING_ENUM_B = 2
	// tags: desc:
	// INT_STRING_ENUM_C = 3
	// tags: desc:[multi ident will skip extract documents and nodes]
	// Multi1 = 1
}

func ExamplePackage_TypeNames() {
	for o := range pkg.TypeNames().Elements() {
		fmt.Println(o.Doc())
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
	// tags:[key1:val_key1_1,val_key1_2][key2:val_key2][key3][key4:val_key4] desc:[IntConstType defines a named constant type with integer underlying in a single `GenDecl`][line1][line2][this is an inline comment]
	// IntConstType
	//
	// tags:[tag1:val1_1,val1_2] desc:[GenDecl defines 2 type, TypeA and TypeB][TypeA doc][line1][line2]
	// TypeA
	//
	// tags:[tag1:val1_1,val1_2] desc:[GenDecl defines 2 type, TypeA and TypeB][TypeB doc][line1][line2]
	// TypeB
	//
	// tags: desc:[IntStringEnum defines a named constant type with integer underlying as an enum type]
	// IntStringEnum
	//
	// tags:[ignore:name] desc:[Structure is a struct type for testing][line1][line2]
	// Structure
	// *Structure.Name: func() string
	// Structure.String: func() string
	// Structure.Value: func() any
	//
	// tags: desc:[StructureAlias is an alias of Structure for testing]
	// StructureAlias
	//
	// tags: desc:[type specs][Int redefines int]
	// Int
	//
	// tags: desc:[type specs][String redefines string]
	// String
	//
	// tags: desc:[type specs][Float alias of float64]
	// Float
	//
}

func ExamplePackage_Functions() {
	for o := range pkg.Functions().Elements() {
		fmt.Println(o.Doc())
		fmt.Printf("%s %s\n\n", o.Ident().Name, o.Type())
	}

	// Output:
	// tags: desc:[Curry function]
	// Curry func() func() int
	//
	// tags: desc:[F a function list call expressions]
	// F func()
}

func ExamplePackages() {
	paths := make([]string, 0)
	fmt.Println("imported in company:")
	for path := range u.Packages {
		if strings.HasPrefix(path, "github.com/xoctopus") {
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
		paths = append(paths, path)
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
	// github.com/xoctopus/pkgx
	// github.com/xoctopus/pkgx/internal
	// github.com/xoctopus/pkgx/testdata
	// github.com/xoctopus/pkgx/testdata/sub
	// github.com/xoctopus/x/contextx
	// github.com/xoctopus/x/misc/must
	// github.com/xoctopus/x/ptrx
	// github.com/xoctopus/x/slicex
	// github.com/xoctopus/x/syncx
	// directs
	// github.com/xoctopus/pkgx
	// github.com/xoctopus/pkgx/internal
	// github.com/xoctopus/pkgx/testdata
	// github.com/xoctopus/pkgx/testdata/sub
	// modules
	// github.com/xoctopus/pkgx
	// github.com/xoctopus/pkgx/testdata
}

func TestWithWorkdir(t *testing.T) {
	dir := filepath.Join(cwd, "..", "x", "reflectx")
	ctx := WithWorkdir(context.Background(), dir)

	x := NewPackages(ctx, "github.com/xoctopus/x/reflectx")
	p := x.Package("github.com/xoctopus/x/reflectx")

	Expect(t, p.SourceDir(), Equal(dir))
	Expect(t, WithWorkdir(ctx, "any"), Equal(ctx))
}
