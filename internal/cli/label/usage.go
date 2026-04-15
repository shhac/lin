package label

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const usageText = `lin label — List Linear issue labels

SUBCOMMANDS:
  label list [--team <team>]  List labels (optionally filtered by team)

OPTIONS:
  --team <team>             Filter by team key or name (e.g., "ENG")
  --limit <n>               Limit results
  --cursor <token>          Pagination cursor for next page

OUTPUT FIELDS:
  list → id, name, color

NOTES:
  Without --team, returns all workspace labels.
  With --team, returns only labels scoped to that team.
  Label names are used as values for --label filters and "issue update labels".`

func registerUsage(label *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Print detailed label command documentation (LLM-optimized)",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(strings.TrimSpace(usageText))
		},
	}
	label.AddCommand(cmd)
}
