package configcmd

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	"github.com/shhac/lin/internal/output"
)

// Register adds the config command group to the parent command.
func Register(parent *cobra.Command) {
	cfg := &cobra.Command{
		Use:   "config",
		Short: "View and update CLI settings",
	}
	output.HandleUnknownCommand(cfg, "Run 'lin config usage' for help")

	registerGet(cfg)
	registerSet(cfg)
	registerReset(cfg)
	registerListKeys(cfg)
	shared.RegisterUsage(cfg, "config", usageText)

	parent.AddCommand(cfg)
}
