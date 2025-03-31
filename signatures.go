package pkgx

import (
	"go/ast"
	"go/types"

	gopkg "golang.org/x/tools/go/packages"
)

type Signature struct {
	node  ast.Node
	ident *ast.Ident
	sig   *types.Signature
}

func (s *Signature) IsZero() bool {
	return s == nil || s.sig == nil
}

func (s *Signature) Name() string {
	if s.ident != nil {
		return s.ident.Name
	}
	return ""
}

func (s *Signature) Node() ast.Node {
	return s.node
}

func (s *Signature) Value() *types.Signature {
	return s.sig
}

func ScanSignatures(p *gopkg.Package) *Set[*Signature] {
	signatures := NewSet[*Signature](p)

	for _, file := range p.Syntax {
		ast.Inspect(file, func(node ast.Node) bool {
			var s = &Signature{node: node}
			switch n := node.(type) {
			case *ast.FuncDecl:
				s.ident = n.Name
				s.sig, _ = p.TypesInfo.TypeOf(n.Name).(*types.Signature)
			case *ast.FuncLit:
				s.sig, _ = p.TypesInfo.TypeOf(n.Type).(*types.Signature)
			}
			if s.IsZero() {
				return true
			}
			signatures.Append(s)
			return true
		})
	}

	for _, file := range p.Syntax {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			sig, ok := p.TypesInfo.TypeOf(call.Fun).(*types.Signature)
			if !ok {
				return true
			}

			scanned := false
			signatures.Range(func(v *Signature) bool {
				if types.Identical(v.sig, sig) {
					scanned = true
					return false
				}
				return true
			})
			if scanned {
				return true
			}
			signatures.Append(&Signature{
				node: node,
				sig:  sig,
			})
			return true
		})
	}

	signatures.Init()
	return signatures
}
