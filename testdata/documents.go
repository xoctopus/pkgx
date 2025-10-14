// Package testdata package level document
//
// comments for testdata package
package testdata

/*
IntConstType defines a named constant type with integer underlying in a single `GenDecl`
line1
line2

+key1=val_key1_1
+key1=val_key1_2
+key2=val_key2
+key3
+key4=
+key4=val_key4
*/
type IntConstType int // this is an inline comment

const (
	// IntConstTypeValue1 doc
	IntConstTypeValue1 IntConstType = iota + 1 // comment 1
	// IntConstTypeValue2 doc
	IntConstTypeValue2 // comment 2
	_                  // placeholder
	IntConstTypeValue3 // comment 3
)

// GenDecl defines 2 type, TypeA and TypeB
type (
	// TypeA doc
	// line1
	// line2
	// +tag1=val1_1
	// +tag1=val1_2
	TypeA int
	// TypeB doc
	// line1
	// line2
	// +tag1=val1_1
	// +tag1=val1_2
	TypeB string
)

// IntStringEnum defines a named constant type with integer underlying as an enum type
type IntStringEnum int

const (
	INT_STRING_ENUM__UNKNOWN IntStringEnum = iota
	// INT_STRING_ENUM_A has doc A
	INT_STRING_ENUM_A
	INT_STRING_ENUM_B // has comment B
	INT_STRING_ENUM_C
)

// multi ident will skip extract documents and nodes
const Multi1, Multi2 = 1, 2

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
