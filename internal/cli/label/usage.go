package label

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const usageText = `lin label — Search, list, and inspect Linear issue labels

SUBCOMMANDS:
  label list [--team <team>] [--name <text>] [--is-group[=false]]   List labels (filterable)
  label search <text> [--team <team>]                               Substring search (case- and accent-insensitive)
  label get <id|name> [--team <team>]                               Single label by UUID or exact name

OPTIONS:
  --team <team>     Filter/disambiguate by team key, name, or UUID
  --name <text>     Exact match (case-insensitive)
  --is-group        Only group labels (--is-group=false for non-groups)
  --limit <n>       Limit results (list, search)
  --cursor <token>  Pagination cursor (list, search)

OUTPUT FIELDS:
  list/search → id, name, color, [description, isGroup, team{id,key,name}, parent{id,name}]
  get         → same fields, single object

NOTES:
  Without --team, labels include the workspace and all teams.
  Workspace-wide labels have no team field. Team-scoped labels include team{key, name}.
  When two labels share a name, "label get" errors and "list --name" shows both with team info.
  Use the resulting label name (with --team) or UUID with "issue new --labels" / "issue update labels".`

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
