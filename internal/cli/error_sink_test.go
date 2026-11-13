package cli

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// PrintError(err.Error()) flattens a structured error to its message: the
// hint is dropped and every failure, a rate limit included, is labelled
// agent-fixable. Route errors through output.HandleGraphQLError (anything
// from a Linear call, resolvers included) or output.WriteError (local errors).
func TestCommandsDoNotFlattenErrors(t *testing.T) {
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(src), "\n") {
			if strings.Contains(line, "PrintError(err.Error())") {
				t.Errorf("%s:%d flattens a structured error; use output.HandleGraphQLError or output.WriteError", path, i+1)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
