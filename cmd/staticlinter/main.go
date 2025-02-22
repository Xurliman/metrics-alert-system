// Package main implements a static analysis tool for Go projects.
//
// StaticLinter integrates various analyzers to enforce best practices and detect potential issues in the codebase.
package main

import (
	"github.com/Xurliman/metrics-alert-system/internal/checkers"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/kisielk/errcheck/errcheck"
)

var (
	// staticCheckAnalyzers defines a list of Staticcheck analyzers that should be included in the analysis.
	staticCheckAnalyzers = []string{"SA4006", "SA5000", "SA6000"}
)

// main is the entry point of the static linter.
//
// It initializes a set of analyzers and runs the multichecker to analyze Go code.
func main() {
	// Define a list of analyzers to be used.
	analyzers := []*analysis.Analyzer{
		printf.Analyzer,            // Checks for incorrect printf-style format strings
		shadow.Analyzer,            // Detects variable shadowing
		structtag.Analyzer,         // Validates struct tags
		checkers.ExitCheckAnalyzer, // Custom analyzer to prevent os.Exit in main files
		errcheck.Analyzer,          // Detects ignored errors in function calls
	}

	// Add selected Staticcheck analyzers based on predefined rules.
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") || slices.Contains(staticCheckAnalyzers, v.Analyzer.Name) {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	// Run the multichecker with the selected analyzers.
	multichecker.Main(analyzers...)
}
