//go:build ignore

// fix-omitempty adds ,omitempty to JSON tags on pointer and slice fields
// in genqlient's generated code. Linear's GraphQL API rejects explicit null
// values in filter/input fields, but genqlient doesn't emit omitempty for
// input types.
package main

import (
	"fmt"
	"os"
	"regexp"
)

var pointerFieldRE = regexp.MustCompile(`(\*\w+)\s+` + "`" + `json:"([^,"]+)"` + "`")
var sliceFieldRE = regexp.MustCompile(`(\[\]\w+)\s+` + "`" + `json:"([^,"]+)"` + "`")

func main() {
	args := os.Args[1:]
	// Skip "--" separator from go run
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: fix-omitempty <file>")
		os.Exit(1)
	}
	path := args[0]
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		os.Exit(1)
	}

	result := pointerFieldRE.ReplaceAll(data, []byte(`${1} `+"`"+`json:"${2},omitempty"`+"`"))
	result = sliceFieldRE.ReplaceAll(result, []byte(`${1} `+"`"+`json:"${2},omitempty"`+"`"))

	if err := os.WriteFile(path, result, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("fix-omitempty: patched %s\n", path)
}
