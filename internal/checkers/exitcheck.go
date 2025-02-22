package checkers

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "checks that os.Exit is not called in main.go files",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Skip if the file is not `cmd/server/main.go` or `cmd/agent/main.go`
		if pass.Pkg.Path() != "cmd/server" && pass.Pkg.Path() != "cmd/agent" {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				// Check if it's a function call to os.Exit
				if ident, ok := call.Fun.(*ast.SelectorExpr); ok {
					if pkgIdent, ok := ident.X.(*ast.Ident); ok && pkgIdent.Name == "os" && ident.Sel.Name == "Exit" {
						pass.Reportf(call.Pos(), "os.Exit should not be used in main.go")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
