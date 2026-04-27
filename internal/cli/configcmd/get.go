package configcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/output"
)

func registerGet(cfg *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get [key]",
		Short: "Show current settings",
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			settings := config.GetSettings()

			if len(args) == 0 {
				output.PrintJSON(settings)
				return
			}

			key := args[0]
			def, ok := settingDefs[key]
			if !ok {
				output.PrintError(fmt.Sprintf("Unknown setting: %s. Valid keys: %s", key, validKeysStr()))
				return
			}

			output.PrintJSON(map[string]any{key: def.get(settings)})
		},
	}
	cfg.AddCommand(cmd)
}
