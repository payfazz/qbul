package qbulcheck

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analizyer = &analysis.Analyzer{
	Name:     "QbulCheck",
	Doc:      "Check qbul usage",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		var obj types.Object
		switch f := call.Fun.(type) {
		case *ast.Ident:
			obj = pass.TypesInfo.ObjectOf(f)
		case *ast.SelectorExpr:
			obj = pass.TypesInfo.ObjectOf(f.Sel)
		default:
			return
		}
		// skip global object
		if obj.Pkg() == nil {
			return
		}
		// skip non method
		if obj.Parent() != nil {
			return
		}
		sig, ok := obj.Type().(*types.Signature)
		if !ok {
			return
		}
		// skip non Builder method
		if sig.Recv().Type().String() != "*github.com/payfazz/qbul.Builder" {
			return
		}
		// only "Add" and "Reset" method
		if name := obj.Name(); name != "Add" && name != "Reset" {
			return
		}
		for _, arg := range call.Args {
			validateArg(pass, arg)
		}
	})
	return nil, nil
}

func validateArg(pass *analysis.Pass, arg ast.Expr) {
	const msg = `must be a string literal or a value with type "Param"`
	if lit, ok := arg.(*ast.BasicLit); ok {
		if lit.Kind != token.STRING {
			pass.ReportRangef(arg, msg)
		}
		return
	}
	if pass.TypesInfo.TypeOf(arg).String() != "github.com/payfazz/qbul.Param" {
		pass.ReportRangef(arg, msg)
	}
}
