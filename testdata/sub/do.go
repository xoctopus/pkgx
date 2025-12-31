package sub

import (
	"context"
)

func Do[Data any, Op interface{ Response() *Data }](ctx context.Context, op Op) (*Data, error) {
	return nil, nil
}

var F = func() int { return 0 }

func Curry() func() string {
	return func() string {
		return ""
	}
}

type AsSel struct{}

type AsSelPtr struct{}

type AsIndex[V any] struct{}

type AsIndexPtr[V any] struct{}

type AsIndexList[V any, K any] struct{}

type AsIndexListPtr[V any, K any] struct{}
