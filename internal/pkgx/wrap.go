package pkgx

import (
	"strings"
	"sync"

	"github.com/xoctopus/x/syncx"
)

const (
	prefix                  = "xwrap_"
	underscore_, underscore = "_", "_u_"
	dot_, dot               = ".", "_d_"
	slash_, slash           = "/", "_s_"
	dash_, dash             = "-", "_k_"
)

func NewWrapper() *Wrapper {
	return sync.OnceValue(func() *Wrapper {
		return &Wrapper{
			p2w: syncx.NewSmap[string, string](),
			w2p: syncx.NewSmap[string, string](),
		}
	})()
}

// Wrapper provides bidirectional mappings between original and wrapped package
// path. To make an identifier with full path can parsed by ast.
// eg:
//
// path/to/pkg.TypeName => path_s_to_s_pkg_d_TypeName
// p2w stores original → wrapped mappings.
// w2p stores wrapped → original mappings.
type Wrapper struct {
	p2w syncx.Map[string, string]
	w2p syncx.Map[string, string]
}

func (w *Wrapper) Clear() {
	w.p2w.Clear()
	w.w2p.Clear()
}

func (w *Wrapper) Unwrap(x string) string {
	if v, ok := w.w2p.Load(x); ok {
		return v
	}
	if _, ok := w.p2w.Load(x); ok {
		return x
	}

	if strings.Contains(x, ".") || strings.Contains(x, "/") {
		return x
	}

	p := x
	p = strings.TrimPrefix(p, prefix)
	p = strings.ReplaceAll(p, slash, slash_)
	p = strings.ReplaceAll(p, dot, dot_)
	p = strings.ReplaceAll(p, dash, dash_)
	p = strings.ReplaceAll(p, underscore, underscore_)

	if !strings.Contains(p, ".") && !strings.Contains(p, "/") {
		x = p
	}

	w.p2w.Store(p, x)
	w.w2p.Store(x, p)

	return p
}

func (w *Wrapper) Wrap(p string) string {
	if x, ok := w.p2w.Load(p); ok {
		return x
	}

	if strings.HasPrefix(p, prefix) {
		return p
	}

	if !strings.Contains(p, ".") && !strings.Contains(p, "/") && !strings.Contains(p, "-") {
		return p
	}

	x := p
	x = strings.ReplaceAll(x, underscore_, underscore)
	x = strings.ReplaceAll(x, dash_, dash)
	x = strings.ReplaceAll(x, dot_, dot)
	x = strings.ReplaceAll(x, slash_, slash)
	x = prefix + x

	w.p2w.Store(p, x)
	w.w2p.Store(x, p)

	return x
}
