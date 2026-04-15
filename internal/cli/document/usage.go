package document

import (
	"fmt"

	"github.com/spf13/cobra"
)

const usageText = `lin document — Document operations (search, list, get, create, update)

SEARCH:
  document search <text>                    Full-text search
    [--include-comments] [--include-archived]

LIST:
  document list                             List documents
    [--project <p>] [--creator <u>] [--include-archived]

GET:
  document get <id>                         Full details + content (markdown)
    Returns: title, content, url, icon, color, project, creator, updatedBy, timestamps

CREATE:
  document new <title>                      Create document
    [--project <p>] [--content <md>] [--icon <emoji>] [--color <hex>]

UPDATE (each is a subcommand):
  document update title <id> <new-title>
  document update content <id> <markdown>
  document update project <id> <project>    Move to project
  document update icon <id> <emoji>
  document update color <id> <hex>

HISTORY:
  document history <id>                       Content edit history (actor IDs + timestamps)

IDS: UUID or slug ID (shown as "slugId" in output). All commands accept either format.
PROJECT: Accepts UUID, slug ID, or name (case-insensitive).
CREATOR: Filter by user ID, name, display name, or email.
CONTENT: Markdown format. Full content returned by "get"; truncated in list/search.
PAGINATION: --limit <n> --cursor <token> on search and list.`

func registerUsage(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Print detailed document command documentation (LLM-optimized)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(usageText)
		},
	}
	parent.AddCommand(cmd)
}
