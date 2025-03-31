package pkgx_test

import (
	"fmt"
	"go/types"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/xoctopus/pkgx"
	_ "github.com/xoctopus/pkgx/testdata"
)

var (
	module   = "github.com/xoctopus/pkgx"
	testdata = "github.com/xoctopus/pkgx/testdata"

	u   *pkgx.Packages
	pkg pkgx.Package
)

func init() {
	u = pkgx.NewPackages(module, testdata)
	pkg = u.Package(testdata)
}

func ExampleScanConstants() {
	constants := pkg.Constants()

	for _, c := range constants {
		fmt.Printf("%s\n", pkg.Doc(c.Node()))
		v := c.Value()
		fmt.Printf("%s [kind:%s] [value:%s]\n\n", c.Name(), v.Kind(), v.String())
	}

	c := pkg.Const("IntConstTypeValue1")
	fmt.Println(c.Value().String())
	c = pkg.ConstByNode(c.Node())
	fmt.Println(c.Value().String())

	// Output:
	// IntConstTypeValue1 tags:[] desc:[IntConstTypeValue1 doc][comment 1]
	// IntConstTypeValue1 [kind:Int] [value:1]
	//
	// IntConstTypeValue2 tags:[] desc:[IntConstTypeValue2 doc][comment 2]
	// IntConstTypeValue2 [kind:Int] [value:2]
	//
	// IntConstTypeValue3 tags:[] desc:[comment 3]
	// IntConstTypeValue3 [kind:Int] [value:3]
	//
	// 1
	// 1
}

func ExampleScanTypenames() {
	typenames := pkg.NamedTypes()

	for _, t := range typenames {
		if doc := pkg.Doc(t.Node()); doc != nil {
			fmt.Printf("%s\n", doc)
		}
		fmt.Printf("%s %s\n", t.Name(), t.Value().String())
		for _, m := range t.Methods() {
			m = t.MethodByName(m.Name())
			fmt.Println(m.Name(), m.Value().String())
		}
		fmt.Println()
	}

	t := pkg.Type("Structure")
	fmt.Println(t.Value().String())
	t = pkg.TypeByNode(t.Node())
	fmt.Println(t.Value().String())

	// Output:
	// IntConstType tags:[key1:val_key1_1,val_key1_2][key2:val_key2][key3:] desc:[IntConstType defines a constant type with integer underlying][line1][line2]
	// IntConstType type github.com/xoctopus/pkgx/testdata.IntConstType int
	//
	// Structure tags:[ignore:name] desc:[Structure is a struct type for testing][line1][line2]
	// Structure type github.com/xoctopus/pkgx/testdata.Structure struct{name string; fieldX any}
	// Name func (*github.com/xoctopus/pkgx/testdata.Structure).Name() string
	// String func (github.com/xoctopus/pkgx/testdata.Structure).String() string
	// Value func (github.com/xoctopus/pkgx/testdata.StructureAlias).Value() any
	//
	// StructureAlias tags:[] desc:[StructureAlias is an alias of Structure for testing]
	// StructureAlias type github.com/xoctopus/pkgx/testdata.StructureAlias = github.com/xoctopus/pkgx/testdata.Structure
	// Name func (*github.com/xoctopus/pkgx/testdata.Structure).Name() string
	// String func (github.com/xoctopus/pkgx/testdata.Structure).String() string
	// Value func (github.com/xoctopus/pkgx/testdata.StructureAlias).Value() any
	//
	// Int tags:[] desc:[type specs][Int redefines int]
	// Int type github.com/xoctopus/pkgx/testdata.Int int
	//
	// String tags:[] desc:[type specs][String redefines string]
	// String type github.com/xoctopus/pkgx/testdata.String string
	//
	// Float tags:[] desc:[type specs][Float alias of float64]
	// Float type github.com/xoctopus/pkgx/testdata.Float = float64
	//
	// with type github.com/xoctopus/pkgx/testdata.with struct{}
	// Call func (github.com/xoctopus/pkgx/testdata.with).Call() (*string, error)
	// With func (github.com/xoctopus/pkgx/testdata.with).With() github.com/xoctopus/pkgx/testdata.with
	//
	// Op type github.com/xoctopus/pkgx/testdata.Op struct{v int}
	// Response func (github.com/xoctopus/pkgx/testdata.Op).Response() *int
	//
	// type github.com/xoctopus/pkgx/testdata.Structure struct{name string; fieldX any}
	// type github.com/xoctopus/pkgx/testdata.Structure struct{name string; fieldX any}
}

func ExampleScanFunctions() {
	functions := pkg.Functions()

	for _, f := range functions {
		if doc := pkg.Doc(f.Node()); doc != nil {
			fmt.Printf("%s\n", doc)
		}
		fmt.Printf("%s %s\n\n", f.Name(), f.Value().String())
	}

	f := pkg.Func("Example")
	fmt.Println(f.Value().String())
	f = pkg.FuncByNode(f.Node())
	fmt.Println(f.Value().String())

	// Output:
	// Example tags:[] desc:[Example a function with nothing return for testing]
	// Example func github.com/xoctopus/pkgx/testdata.Example()
	//
	// FuncSingleReturn tags:[] desc:[FuncSingleReturn a function with single return for testing]
	// FuncSingleReturn func github.com/xoctopus/pkgx/testdata.FuncSingleReturn() any
	//
	// FuncSelectExprReturn tags:[] desc:[FuncSelectExprReturn a function returns a struct field for testing]
	// FuncSelectExprReturn func github.com/xoctopus/pkgx/testdata.FuncSelectExprReturn() string
	//
	// FuncWithCall tags:[] desc:[FuncWithCall a function returns other function result and type assert]
	// FuncWithCall func github.com/xoctopus/pkgx/testdata.FuncWithCall() (any, github.com/xoctopus/pkgx/testdata.String)
	//
	// FuncReturnInterfaceCallMulti func github.com/xoctopus/pkgx/testdata.FuncReturnInterfaceCallMulti() (any, error)
	//
	// FuncReturnInterfaceCallSingle func github.com/xoctopus/pkgx/testdata.FuncReturnInterfaceCallSingle() any
	//
	// FuncReturnsNamedValue func github.com/xoctopus/pkgx/testdata.FuncReturnsNamedValue() (a any, b github.com/xoctopus/pkgx/testdata.String)
	//
	// FuncReturnsNamedValueAndOtherFunc func github.com/xoctopus/pkgx/testdata.FuncReturnsNamedValueAndOtherFunc() (a any, b github.com/xoctopus/pkgx/testdata.String, err error)
	//
	// FuncReturnsInSwitch func github.com/xoctopus/pkgx/testdata.FuncReturnsInSwitch(v string) (a any, b github.com/xoctopus/pkgx/testdata.String)
	//
	// FuncReturnsInIf func github.com/xoctopus/pkgx/testdata.FuncReturnsInIf(v string) (a any, b github.com/xoctopus/pkgx/testdata.String)
	//
	// FuncCallWithFuncLit func github.com/xoctopus/pkgx/testdata.FuncCallWithFuncLit() (a any, b github.com/xoctopus/pkgx/testdata.String)
	//
	// With func github.com/xoctopus/pkgx/testdata.With() github.com/xoctopus/pkgx/testdata.with
	//
	// FuncWithCallChain func github.com/xoctopus/pkgx/testdata.FuncWithCallChain() (any, error)
	//
	// FuncWithSub func github.com/xoctopus/pkgx/testdata.FuncWithSub() (any, error)
	//
	// curry func github.com/xoctopus/pkgx/testdata.curry(b bool) func() int
	//
	// FuncCurryCall func github.com/xoctopus/pkgx/testdata.FuncCurryCall() any
	//
	// func github.com/xoctopus/pkgx/testdata.Example()
	// func github.com/xoctopus/pkgx/testdata.Example()
}

func TestFunction_CanRefByValue(t *testing.T) {

	f := pkg.Func("Example")
	NewWithT(t).Expect(f.CanRefByValue()).To(BeFalse())

	s := pkg.Type("Structure")
	NewWithT(t).Expect(len(s.Methods())).To(Equal(3))

	NewWithT(t).Expect(s.MethodByName("Name").CanRefByValue()).To(BeFalse())
	NewWithT(t).Expect(s.MethodByName("String").CanRefByValue()).To(BeTrue())
}

func ExampleScanSignatures() {
	_, filename, _, _ := runtime.Caller(0)
	cwd := filepath.Dir(filename)

	signatures := pkg.Signatures()
	for _, s := range signatures {
		node := s.Node()
		if name := s.Name(); name == "" {
			fmt.Printf("%T %s\n", node, s.Value())
		} else {
			fmt.Printf("%s %T %s\n", name, node, s.Value())
		}
		fmt.Printf("%s\n", pkgx.PositionOf(pkg.GoPackage(), node).WithCwd(cwd).String())
	}

	// Output:
	// Name *ast.FuncDecl func() string
	// testdata/testdata.go:44:1
	// testdata/testdata.go:46:2
	// String *ast.FuncDecl func() string
	// testdata/testdata.go:48:1
	// testdata/testdata.go:50:2
	// Value *ast.FuncDecl func() any
	// testdata/testdata.go:52:1
	// testdata/testdata.go:54:2
	// Example *ast.FuncDecl func()
	// testdata/testdata.go:67:1
	// testdata/testdata.go:67:18
	// FuncSingleReturn *ast.FuncDecl func() any
	// testdata/testdata.go:70:1
	// testdata/testdata.go:80:2
	// *ast.FuncLit func() any
	// testdata/testdata.go:71:6
	// testdata/testdata.go:74:3
	// FuncSelectExprReturn *ast.FuncDecl func() string
	// testdata/testdata.go:83:1
	// testdata/testdata.go:87:2
	// FuncWithCall *ast.FuncDecl func() (any, github.com/xoctopus/pkgx/testdata.String)
	// testdata/testdata.go:90:1
	// testdata/testdata.go:92:2
	// FuncReturnInterfaceCallMulti *ast.FuncDecl func() (any, error)
	// testdata/testdata.go:94:1
	// testdata/testdata.go:96:2
	// *ast.CallExpr func(p []byte) (n int, err error)
	// testdata/testdata.go:95:9
	// testdata/testdata.go:95:34
	// FuncReturnInterfaceCallSingle *ast.FuncDecl func() any
	// testdata/testdata.go:98:1
	// testdata/testdata.go:100:2
	// *ast.CallExpr func() error
	// testdata/testdata.go:99:9
	// testdata/testdata.go:99:31
	// FuncReturnsNamedValue *ast.FuncDecl func() (a any, b github.com/xoctopus/pkgx/testdata.String)
	// testdata/testdata.go:102:1
	// testdata/testdata.go:105:2
	// FuncReturnsNamedValueAndOtherFunc *ast.FuncDecl func() (a any, b github.com/xoctopus/pkgx/testdata.String, err error)
	// testdata/testdata.go:107:1
	// testdata/testdata.go:111:2
	// *ast.CallExpr func(message string) error
	// testdata/testdata.go:110:15
	// testdata/testdata.go:110:32
	// FuncReturnsInSwitch *ast.FuncDecl func(v string) (a any, b github.com/xoctopus/pkgx/testdata.String)
	// testdata/testdata.go:113:1
	// testdata/testdata.go:128:2
	// FuncReturnsInIf *ast.FuncDecl func(v string) (a any, b github.com/xoctopus/pkgx/testdata.String)
	// testdata/testdata.go:130:1
	// testdata/testdata.go:146:2
	// FuncCallWithFuncLit *ast.FuncDecl func() (a any, b github.com/xoctopus/pkgx/testdata.String)
	// testdata/testdata.go:148:1
	// testdata/testdata.go:153:2
	// *ast.FuncLit func() any
	// testdata/testdata.go:149:10
	// testdata/testdata.go:151:3
	// With *ast.FuncDecl func() github.com/xoctopus/pkgx/testdata.with
	// testdata/testdata.go:157:1
	// testdata/testdata.go:159:2
	// With *ast.FuncDecl func() github.com/xoctopus/pkgx/testdata.with
	// testdata/testdata.go:161:1
	// testdata/testdata.go:163:2
	// Call *ast.FuncDecl func() (*string, error)
	// testdata/testdata.go:165:1
	// testdata/testdata.go:167:2
	// *ast.CallExpr func(v string) *string
	// testdata/testdata.go:166:9
	// testdata/testdata.go:166:21
	// FuncWithCallChain *ast.FuncDecl func() (any, error)
	// testdata/testdata.go:169:1
	// testdata/testdata.go:171:2
	// Response *ast.FuncDecl func() *int
	// testdata/testdata.go:177:1
	// testdata/testdata.go:179:2
	// FuncWithSub *ast.FuncDecl func() (any, error)
	// testdata/testdata.go:181:1
	// testdata/testdata.go:183:2
	// *ast.CallExpr func(ctx context.Context, op *github.com/xoctopus/pkgx/testdata.Op) (*int, error)
	// testdata/testdata.go:182:9
	// testdata/testdata.go:182:44
	// *ast.CallExpr func() context.Context
	// testdata/testdata.go:182:16
	// testdata/testdata.go:182:36
	// curry *ast.FuncDecl func(b bool) func() int
	// testdata/testdata.go:185:1
	// testdata/testdata.go:198:2
	// *ast.FuncLit func() func() int
	// testdata/testdata.go:187:10
	// testdata/testdata.go:191:4
	// *ast.FuncLit func() int
	// testdata/testdata.go:188:11
	// testdata/testdata.go:190:5
	// *ast.FuncLit func() func() int
	// testdata/testdata.go:193:9
	// testdata/testdata.go:197:3
	// *ast.FuncLit func() int
	// testdata/testdata.go:194:10
	// testdata/testdata.go:196:4
	// FuncCurryCall *ast.FuncDecl func() any
	// testdata/testdata.go:200:1
	// testdata/testdata.go:202:2
}

func TestXpkg(t *testing.T) {
	t.Run("Unwrap", func(t *testing.T) {
		NewWithT(t).Expect(pkg.Unwrap().Path()).To(Equal(testdata))
	})
	t.Run("GoPackage", func(t *testing.T) {
		NewWithT(t).Expect(pkg.GoPackage().PkgPath).To(Equal(testdata))
	})
	t.Run("Package", func(t *testing.T) {
		NewWithT(t).Expect(pkg.Package("")).To(BeNil())
		NewWithT(t).Expect(pkg.Package("github.com/pkg/errors")).NotTo(BeNil())
		NewWithT(t).Expect(pkg.Package("github.com/pkg/errors")).NotTo(BeNil())
	})

	t.Run("SourceDir", func(t *testing.T) {
		_, filename, _, _ := runtime.Caller(0)
		cwd := filepath.Dir(filename)

		dir, _ := filepath.Rel(cwd, pkg.SourceDir())
		NewWithT(t).Expect(dir).To(Equal("testdata"))

		NewWithT(t).Expect(u.Package("io").SourceDir()).To(Equal(""))

		dir, _ = filepath.Rel(cwd, u.Package(module).SourceDir())
		NewWithT(t).Expect(dir).To(Equal("."))

		dir, _ = filepath.Rel(cwd, pkg.SourceDir())
		NewWithT(t).Expect(dir).To(Equal("testdata"))

		p := u.Package("github.com/xoctopus/pkgx/testdata/sub")
		dir, _ = filepath.Rel(cwd, p.SourceDir())
		NewWithT(t).Expect(dir).To(Equal("testdata/sub"))
	})

	t.Run("Module", func(t *testing.T) {
		NewWithT(t).Expect(pkg.Module().Path).To(Equal(testdata))
	})

	t.Run("Files", func(t *testing.T) {
		files := pkg.Files()
		names := make([]string, 0, len(files))
		for _, f := range files {
			names = append(names, f.Name.Name)
		}
		NewWithT(t).Expect(names).To(ConsistOf("testdata"))
	})

	obj := pkg.Func("Example")
	filename := pkgx.PositionOf(pkg.GoPackage(), obj.Node()).Filename()

	t.Run("FileSet", func(t *testing.T) {
		fun := pkg.FileSet().File(obj.Node().Pos())
		NewWithT(t).Expect(fun.Name()).To(Equal(filename))
	})

	t.Run("Position", func(t *testing.T) {
		pos := pkg.Position(obj.Node().Pos())
		NewWithT(t).Expect(pos.Filename).To(Equal(filename))
	})

	t.Run("ObjectOf", func(t *testing.T) {
		NewWithT(t).Expect(pkg.ObjectOf(obj.Ident())).To(Equal(obj.Value()))
	})

	t.Run("Eval", func(t *testing.T) {
		tv, err := pkg.Eval(obj.Ident())
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(types.Identical(tv.Type, obj.Value().Signature())).To(BeTrue())

		_, err = pkg.Eval(nil)
		NewWithT(t).Expect(err).NotTo(BeNil())
	})
}

func TestIdentNode(t *testing.T) {
	ident := pkg.Type("Structure").MethodByName("Name").Ident()
	NewWithT(t).Expect(pkgx.IdentNode(ident)).To(BeNil())
}
