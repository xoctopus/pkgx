package pkgx

import (
	"go/constant"
	"go/types"
)

type Constant struct {
	Object[*types.Const]
}

func (c *Constant) Value() constant.Value {
	return c.Exposer().Val()
}
