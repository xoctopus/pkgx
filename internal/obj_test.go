package internal_test

import (
	"go/ast"
	"go/types"
	"testing"

	"github.com/xoctopus/pkgx"
	. "github.com/xoctopus/pkgx/internal"
	. "github.com/xoctopus/x/testx"
)

var (
	// sub      = "github.com/xoctopus/pkgx/testdata/sub"
	// cwd, _   = os.Getwd()
	u   = pkgx.NewPackages("github.com/xoctopus/pkgx", "github.com/xoctopus/pkgx/testdata")
	pkg = u.Package("github.com/xoctopus/pkgx/testdata")
)

func TestNewObject(t *testing.T) {
	t.Run("InvalidObject", func(t *testing.T) {
		o := NewObject[*types.Func](nil, nil, nil, nil)
		Expect(t, o.IsNil(), BeTrue())
		Expect(t, o.Node(), BeNil[ast.Node]())
		Expect(t, o.Name(), HaveLen[string](0))
		Expect(t, o.Ident(), BeNil[*ast.Ident]())
		Expect(t, o.Underlying(), BeNil[*types.Func]())
		Expect(t, o.Doc(), BeNil[*Doc]())
		Expect(t, o.Type(), BeNil[types.Type]())
		Expect(t, o.TypeName(), HaveLen[string](0))
	})
}

func TestObject(t *testing.T) {
	constants := pkg.Constants()
	typenames := pkg.TypeNames()
	functions := pkg.Functions()

	c := constants.ElementByName("IntConstTypeValue1")
	f := functions.ElementByName("F")
	n := typenames.ElementByName("Structure")

	Expect(t, c.Value().String(), Equal("1"))
	Expect(t, f.PtrRecv(), BeFalse())

	Expect(t, c.IsNil(), BeFalse())
	Expect(t, c.Name(),
		Equal("IntConstTypeValue1"),
		Equal(c.Ident().Name),
	)

	Expect(t, c.Type(), NotBeNil[types.Type]())
	Expect(t, f.Type(), NotBeNil[types.Type]())

	Expect(t, c.TypeName(), Equal("IntConstType"))
	Expect(t, f.TypeName(), Equal(""))
	Expect(t, n.TypeName(), Equal("Structure"))

	Expect(t, functions.Len(), Equal(2))

	var node ast.Node
	for node = range functions.Nodes() {
		if node.Pos() == f.Node().Pos() && node.End() == f.Node().End() {
			x, ok := node.(*ast.FuncDecl)
			Expect(t, ok, BeTrue())
			Expect(t, x.Name, Equal(f.Ident()))
			break
		}
	}

	functions.Add(&Function{Object: NewObject(nil, nil, &types.Func{}, nil)})
	Expect(t, functions.Len(), Equal(2))

	fu := functions.Underlying(node)
	Expect(t, fu, Equal[*types.Func](f.Underlying()))
	fu = functions.Underlying(Node{1, 1})
	Expect(t, fu, BeNil[*types.Func]())

	for fu = range functions.Underlyings() {
		if fu.Name() == "F" {
			Expect(t, fu, Equal(f.Underlying()))
			break
		}
	}

	fo := functions.Element(node)
	Expect(t, fo, Equal(f))

	for fo = range functions.Elements() {
		if fo.Name() == "F" {
			Expect(t, fo, Equal(f))
			break
		}
	}

	n.Methods().Range(func(key string, f *Function) bool {
		Expect(t, f.Name(), Equal(key))
		if key == "Name" {
			Expect(t, f.PtrRecv(), BeTrue())
		}
		Expect(t, n.Method(key), Equal(f))
		return true
	})
}
