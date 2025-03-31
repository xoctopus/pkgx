package sub

// Structure is a struct type for testing in testdata/sub
// line1
type Structure struct {
}

func (v *Structure) Name() string {
	return "name"
}

func (v Structure) String() string {
	return "structure"
}
