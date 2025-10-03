package internal

import (
	"go/types"
)

type Function struct{ Object[*types.Func] }

func (f *Function) PtrRecv() bool {
	recv := f.Underlying().Signature().Recv()
	if recv == nil {
		return false
	}
	_, ok := recv.Type().(*types.Pointer)
	return ok
}
