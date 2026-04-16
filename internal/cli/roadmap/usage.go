package roadmap

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const usageText = `lin roadmap — Browse Linear roadmaps (read-only)

LIST:
  roadmap list                            List all roadmaps
    [--limit] [--cursor]
    Returns per item: id, slugId, url, name, description, owner

GET:
  roadmap get <id>             Roadmap summary: id, slugId, url, name, description,
                               owner, creator, createdAt
  roadmap projects <id>        Projects linked to a roadmap
    [--limit] [--cursor]
    Returns per item: id, slugId, url, name, status, progress,
    lead, startDate, targetDate

IDS: <id> accepts UUID, slug ID, or roadmap name.
PAGINATION: --limit <n> --cursor <token> on list and projects.

NOTE: Roadmaps are read-only. Use "project" commands to modify projects in a roadmap.`

func registerUsage(roadmap *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Print detailed roadmap command documentation (LLM-optimized)",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(strings.TrimSpace(usageText))
		},
	}
	roadmap.AddCommand(cmd)
}
