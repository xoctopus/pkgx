package testdata

import (
	"fmt"

	"github.com/xoctopus/pkgx/testdata/sub"
)

// function var
var ff = func() int { return 0 }

// Curry function
func Curry() func() int {
	return func() int {
		return func() int {
			return 0
		}()
	}
}

// F a function list call expressions
func F() {
	// call func val from current package
	ff()

	// call func literal by name
	f := func() {}
	f()

	// call func literal directly
	_ = func() int { return 1 }()

	// call method by a value from current package
	_ = (&Structure{}).Name()
	_ = new(Structure).Name()
	s := new(Structure)
	_ = s.Name()

	// // call func val from other package
	_ = sub.F()

	// chain call
	_ = (sub.Structure{}).With().Name()

	// curry call
	Curry()()
	sub.Curry()()

	_ = fmt.Sprintf("")
}
