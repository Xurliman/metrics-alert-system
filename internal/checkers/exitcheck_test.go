package checkers

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestExitCheckAnalyzer verifies that the ExitCheckAnalyzer detects os.Exit calls in main.go.
func TestExitCheckAnalyzer(t *testing.T) {
	// Define the testdata directory containing test Go files.
	testdata := analysistest.TestData()

	// Run the analyzer on the testdata directory.
	analysistest.Run(t, testdata, ExitCheckAnalyzer, "cmd/server", "cmd/agent")
}
