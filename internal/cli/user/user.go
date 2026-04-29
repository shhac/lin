package user

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	"github.com/shhac/lin/internal/output"
)

// Register adds the user command group to the parent command.
func Register(parent *cobra.Command) {
	user := &cobra.Command{
		Use:   "user",
		Short: "User operations",
	}
	output.HandleUnknownCommand(user, "Run 'lin user usage' for help")

	registerList(user)
	registerMe(user)
	registerSearch(user)
	shared.RegisterUsage(user, "user", usageText)

	parent.AddCommand(user)
}
