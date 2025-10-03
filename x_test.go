package pkgx_test

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"sort"
	"testing"

	. "github.com/xoctopus/pkgx"
	_ "github.com/xoctopus/pkgx/testdata"
	"github.com/xoctopus/x/mapx"
	. "github.com/xoctopus/x/testx"
)

var (
	module   = "github.com/xoctopus/pkgx"
	testdata = "github.com/xoctopus/pkgx/testdata"
	sub      = "github.com/xoctopus/pkgx/testdata/sub"
	cwd, _   = os.Getwd()

	u   = NewPackages(module, testdata)
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
}

func ExamplePackage_constants() {
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

func ExamplePackage_typenames() {
	for o := range pkg.TypeNames().Elements() {
		fmt.Println(o.Doc())
		fmt.Println(o.Ident().Name)
		methods := mapx.Keys(o.Methods())
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

func ExamplePackage_functions() {
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

var (
	_ ast.Node = &ast.ArrayType{}
	_ ast.Node = &ast.AssignStmt{}
	_ ast.Node = &ast.BadDecl{}
	_ ast.Node = &ast.BadExpr{}
	_ ast.Node = &ast.BadStmt{}
	_ ast.Node = &ast.BasicLit{}
	_ ast.Node = &ast.BinaryExpr{}
	_ ast.Node = &ast.BranchStmt{}
	_ ast.Node = &ast.CallExpr{}
	_ ast.Node = &ast.CaseClause{}
	_ ast.Node = &ast.ChanType{}
	_ ast.Node = &ast.CommClause{}
	_ ast.Node = &ast.Comment{}
	_ ast.Node = &ast.CommentGroup{}
	_ ast.Node = &ast.CompositeLit{}
	_ ast.Node = ast.Decl(nil)
	_ ast.Node = &ast.DeclStmt{}
	_ ast.Node = &ast.DeferStmt{}
	_ ast.Node = &ast.Ellipsis{}
	_ ast.Node = &ast.EmptyStmt{}
	_ ast.Node = ast.Expr(nil)
	_ ast.Node = &ast.ExprStmt{}
	_ ast.Node = &ast.Field{}
	_ ast.Node = &ast.FieldList{}
	_ ast.Node = &ast.File{}
	_ ast.Node = &ast.ForStmt{}
	_ ast.Node = &ast.FuncDecl{}
	_ ast.Node = &ast.FuncLit{}
	_ ast.Node = &ast.FuncType{}
	_ ast.Node = &ast.GenDecl{}
	_ ast.Node = &ast.GoStmt{}
	_ ast.Node = &ast.Ident{}
	_ ast.Node = &ast.IfStmt{}
	_ ast.Node = &ast.ImportSpec{}
	_ ast.Node = &ast.IncDecStmt{}
	_ ast.Node = &ast.IndexExpr{}
	_ ast.Node = &ast.IndexListExpr{}
	_ ast.Node = &ast.InterfaceType{}
	_ ast.Node = &ast.KeyValueExpr{}
	_ ast.Node = &ast.LabeledStmt{}
	_ ast.Node = &ast.MapType{}
	// _ ast.Node = &ast.Package{}
	_ ast.Node = &ast.ParenExpr{}
	_ ast.Node = &ast.RangeStmt{}
	_ ast.Node = &ast.ReturnStmt{}
	_ ast.Node = &ast.SelectStmt{}
	_ ast.Node = &ast.SelectorExpr{}
	_ ast.Node = &ast.SendStmt{}
	_ ast.Node = &ast.SliceExpr{}
	_ ast.Node = ast.Spec(nil)
	_ ast.Node = &ast.StarExpr{}
	_ ast.Node = ast.Stmt(nil)
	_ ast.Node = &ast.StructType{}
	_ ast.Node = &ast.SwitchStmt{}
	_ ast.Node = &ast.TypeAssertExpr{}
	_ ast.Node = &ast.TypeSpec{}
	_ ast.Node = &ast.TypeSwitchStmt{}
	_ ast.Node = &ast.UnaryExpr{}
	_ ast.Node = &ast.ValueSpec{}
)
