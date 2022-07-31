//go:build !e2e
// +build !e2e

package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestE2E(t *testing.T) {
	currentPath, _ := os.Getwd()
	workspacePath, _ := filepath.Abs(currentPath + "/..")

	passed := 0
	failed := 0

	exampleList, _ := filepath.Glob(fmt.Sprintf("%s/example/*", workspacePath))

	for _, dir := range exampleList {
		t, f := RunExampleTest(workspacePath, dir)
		passed += t
		failed += f
	}
	fmt.Printf("[E2E Test]: %d passed, %d failed\n", passed, failed)
	if failed != 0 {
		t.Fail()
	}
}
