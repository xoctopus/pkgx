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
