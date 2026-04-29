package label

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	"github.com/shhac/lin/internal/output"
)

// Register adds the label command group to the parent command.
func Register(parent *cobra.Command) {
	label := &cobra.Command{
		Use:   "label",
		Short: "Label operations",
	}
	output.HandleUnknownCommand(label, "Run 'lin label usage' for help")

	registerList(label)
	registerSearch(label)
	registerGet(label)
	shared.RegisterUsage(label, "label", usageText)

	parent.AddCommand(label)
}
