// Package testdata package level document
//
// comments for testdata package
package testdata

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/xoctopus/x/ptrx"

	"github.com/xoctopus/pkgx/testdata/sub"
)

// IntConstType defines a constant type with integer underlying
// line1
// line2
// +key1=val_key1_1
// +key1=val_key1_2
// +key2=val_key2
// +key3
type IntConstType int

const (
	// IntConstTypeValue1 doc
	IntConstTypeValue1 IntConstType = iota + 1 // comment 1
	// IntConstTypeValue2 doc
	IntConstTypeValue2 // comment 2
	IntConstTypeValue3 // comment 3
)

// Structure is a struct type for testing
// line1
// line2
// +ignore=name
type Structure struct {
	name   string
	fieldX any
}

// StructureAlias is an alias of Structure for testing
type StructureAlias = Structure

func (v *Structure) Name() string {
	return "name"
}

func (v Structure) String() string {
	return "structure"
}

func (v StructureAlias) Value() any {
	return 1
}

// type specs
type (
	// Int redefines int
	Int int
	// String redefines string
	String string
	// Float alias of float64
	Float = float64
)

// Example a function with nothing return for testing
func Example() {}

// FuncSingleReturn a function with single return for testing
func FuncSingleReturn() any {
	_ = func() any {
		v := false
		return !v
	}

	var v any
	v = "1"
	v = 2
	return v
}

// FuncSelectExprReturn a function returns a struct field for testing
func FuncSelectExprReturn() string {
	v := struct{ s string }{}
	v.s = "1"
	return v.s
}

// FuncWithCall a function returns other function result and type assert
func FuncWithCall() (any, String) {
	return FuncSingleReturn(), String(FuncSelectExprReturn())
}

func FuncReturnInterfaceCallMulti() (any, error) {
	return io.Writer(nil).Write(nil)
}

func FuncReturnInterfaceCallSingle() any {
	return io.Closer(nil).Close()
}

func FuncReturnsNamedValue() (a any, b String) {
	a, b = FuncWithCall()
	return
}

func FuncReturnsNamedValueAndOtherFunc() (a any, b String, err error) {
	a = "a"
	b = "b"
	return a, b, errors.New("any")
}

func FuncReturnsInSwitch(v string) (a any, b String) {
	switch v {
	case "1":
		a = "a1"
		b = "b1"
		return
	case "2":
		a = "a2"
		b = "b2"
		return
	default:
		a = "a3"
		b = "b3"
		return
	}
}

func FuncReturnsInIf(v string) (a any, b String) {
	if v == "1" {
		a = "a1"
		b = "b1"
		return
	} else if v == "2" {
		a = "a2"
		b = "b2"
		return
	} else if true {
		a = "a3"
		b = "b3"
		return
	}

	return FuncWithCall()
}

func FuncCallWithFuncLit() (a any, b String) {
	call := func() any {
		return 1
	}
	return call(), "s"
}

type with struct{}

func With() with {
	return with{}
}

func (w with) With() with {
	return w
}

func (w with) Call() (*string, error) {
	return ptrx.Ptr(""), nil
}

func FuncWithCallChain() (any, error) {
	return With().With().Call()
}

type Op struct {
	v int
}

func (o Op) Response() *int {
	return &o.v
}

func FuncWithSub() (any, error) {
	return sub.Do(context.Background(), &Op{})
}

func curry(b bool) func() int {
	if b {
		return func() func() int {
			return func() int {
				return 1
			}
		}()
	}
	return func() func() int {
		return func() int {
			return 0
		}
	}()
}

func FuncCurryCall() any {
	return curry(true)
}
