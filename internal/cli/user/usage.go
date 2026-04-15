package user

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const usageText = `lin user — Look up Linear workspace users

SUBCOMMANDS:
  user search <text>         Search users by name, email, or display name
  user list [--team <team>]  List users (optionally filtered by team)
  user me                    Current authenticated user + organization

OPTIONS (search):
  --limit <n>               Limit results
  --cursor <token>          Pagination cursor for next page

OPTIONS (list):
  --team <team>             Filter by team key or name (e.g., "ENG")
  --limit <n>               Limit results
  --cursor <token>          Pagination cursor for next page

OUTPUT FIELDS:
  search → id, name, email, displayName
  list   → id, name, email, displayName
  me     → id, name, email, displayName, organization (id, name)

NOTES:
  User IDs are UUIDs. Use "user me" to get the current user's ID for
  operations like --assignee filtering.
  When --team is specified, returns only members of that team.
  Search matches against name, email, and displayName (case-insensitive).`

func registerUsage(user *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Print detailed user command documentation (LLM-optimized)",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(strings.TrimSpace(usageText))
		},
	}
	user.AddCommand(cmd)
}
