package auth

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output"
)

// Register adds the auth command group to the parent command.
func Register(parent *cobra.Command) {
	auth := &cobra.Command{
		Use:   "auth",
		Short: "Authentication management",
	}
	output.HandleUnknownCommand(auth, "Run 'lin auth usage' for help")

	registerLogin(auth)
	registerStatus(auth)
	registerLogout(auth)
	registerWorkspace(auth)
	registerUsage(auth)

	parent.AddCommand(auth)
}
