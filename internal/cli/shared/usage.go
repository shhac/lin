// Package shared holds helpers used by every lin subcommand.
package shared

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// RegisterUsage wires the canonical `<verb> usage` subcommand. Centralised so
// the per-package fmt+strings imports and 4-line cobra block aren't duplicated.
func RegisterUsage(parent *cobra.Command, verb, text string) {
	parent.AddCommand(&cobra.Command{
		Use:   "usage",
		Short: "Show detailed reference for " + verb,
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(strings.TrimSpace(text))
		},
	})
}
