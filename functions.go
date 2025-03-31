package pkgx

import (
	"go/ast"
	"go/types"

	gopkg "golang.org/x/tools/go/packages"
)

type Function struct {
	node     ast.Node
	ident    *ast.Ident
	function *types.Func
}

func (f *Function) IsZero() bool {
	return f == nil || f.function == nil
}

func (f *Function) Name() string {
	return f.ident.Name
}

func (f *Function) Ident() *ast.Ident {
	return f.ident
}

func (f *Function) Node() ast.Node {
	return f.node
}

func (f *Function) Value() *types.Func {
	return f.function
}

func (f *Function) CanRefByValue() bool {
	recv := f.function.Signature().Recv()
	if recv == nil {
		return false
	}
	if _, ptr := recv.Type().(*types.Pointer); ptr {
		return false
	}
	return true
}

func ScanFunctions(pkg *gopkg.Package) *Set[*Function] {
	functions := NewSet[*Function](pkg)

	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(node ast.Node) bool {
			f, ok := node.(*ast.FuncDecl)
			if !ok || f.Recv != nil {
				return true
			}
			functions.Append(&Function{
				node:     node,
				ident:    f.Name,
				function: pkg.TypesInfo.ObjectOf(f.Name).(*types.Func),
			})
			return true
		})
	}

	functions.Init()
	return functions
}
