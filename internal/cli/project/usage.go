package project

import (
	"fmt"

	"github.com/spf13/cobra"
)

const usageText = `lin project — Project operations (search, list, get, create, update)

SEARCH & LIST:
  project search <text>                   Full-text search for projects
    [--limit] [--cursor]
  project list                            List all projects
    [--team <team>] [--status <status>] [--limit] [--cursor]

GET:
  project get <id>             Project summary: id, slugId, url, name, description, content,
                               status, progress, lead, startDate, targetDate, milestones[]
  project issues <id>          Issues in a project
    [--status] [--assignee] [--priority] [--limit] [--cursor]

CREATE:
  project new <name> --team <teams>       --team required (comma-separated for multi-team)
    [--description <desc>] [--lead <user>] [--start-date <YYYY-MM-DD>]
    [--target-date <YYYY-MM-DD>] [--status <status>] [--content <markdown>]

UPDATE (each is a subcommand):
  project update title <id> <new-title>
  project update status <id> <new-status>
  project update description <id> <description>
  project update content <id> <markdown>
  project update lead <id> <user>
  project update start-date <id> <YYYY-MM-DD>
  project update target-date <id> <YYYY-MM-DD>
  project update priority <id> <priority>    none|urgent|high|medium|low
  project update icon <id> <emoji>
  project update color <id> <hex>
  project update labels <id> <labels>        Replace project labels (comma-separated names or UUIDs).
                                             Resolved against project labels only — see
                                             "lin label list --type project". Replace semantics: any
                                             previously-set label not listed is removed.

LIFECYCLE:
  project delete <id>          Delete (trash) a project
  project unarchive <id>       Restore a trashed or archived project

IDS: <id> accepts UUID, slug ID, or project name.
TEAM: Team key (ENG), name, or UUID.
LEAD: Name, email, or user ID.
PROJECT STATUS: backlog|planned|started|paused|completed|canceled
PAGINATION: --limit <n> --cursor <token> on search, list, and issues.`

func registerUsage(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Print detailed project command documentation (LLM-optimized)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(usageText)
		},
	}
	parent.AddCommand(cmd)
}
