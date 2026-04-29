package team

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	"github.com/shhac/lin/internal/output"
)

// Register adds the team command group to the parent command.
func Register(parent *cobra.Command) {
	team := &cobra.Command{
		Use:   "team",
		Short: "Team operations",
	}
	output.HandleUnknownCommand(team, "Run 'lin team usage' for help")

	registerList(team)
	registerGet(team)
	registerStates(team)
	shared.RegisterUsage(team, "team", usageText)

	parent.AddCommand(team)
}
