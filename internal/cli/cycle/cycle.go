package cycle

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output"
)

// Register adds the cycle command group to the parent command.
func Register(parent *cobra.Command) {
	cycle := &cobra.Command{
		Use:   "cycle",
		Short: "Cycle operations",
	}
	output.HandleUnknownCommand(cycle, "Run 'lin cycle usage' for help")

	registerList(cycle)
	registerGet(cycle)
	registerUsage(cycle)

	parent.AddCommand(cycle)
}
