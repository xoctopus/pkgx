package sub

import (
	"context"
)

func Do[Data any, Op interface{ Response() *Data }](ctx context.Context, op Op) (*Data, error) {
	return nil, nil
}
