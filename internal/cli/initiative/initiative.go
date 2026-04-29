package initiative

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	"github.com/shhac/lin/internal/output"
)

// Register adds the initiative command group to the parent command.
func Register(parent *cobra.Command) {
	initiative := &cobra.Command{
		Use:   "initiative",
		Short: "Initiative operations",
	}
	parent.AddCommand(initiative)

	registerSearch(initiative)
	registerList(initiative)
	registerGet(initiative)
	registerProjects(initiative)
	registerNew(initiative)
	registerUpdate(initiative)
	registerArchive(initiative)
	registerUnarchive(initiative)
	registerDelete(initiative)
	shared.RegisterUsage(initiative, "initiative", usageText)
	output.HandleUnknownCommand(initiative, "Run 'lin initiative usage' for help")
}
