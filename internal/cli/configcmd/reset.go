package configcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/output"
)

func registerReset(cfg *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "reset [key]",
		Short: "Reset settings to defaults",
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				if err := config.ResetSettings(); err != nil {
					output.PrintError(err.Error())
					return
				}
				output.PrintJSON(map[string]any{"reset": "all"})
				return
			}

			key := args[0]
			def, ok := settingDefs[key]
			if !ok {
				output.PrintError(fmt.Sprintf("Unknown setting: %s. Valid keys: %s", key, validKeysStr()))
				return
			}

			settings := config.GetSettings()
			def.reset(settings)

			if err := config.UpdateSettings(settings); err != nil {
				output.PrintError(err.Error())
				return
			}
			output.PrintJSON(map[string]any{"reset": key})
		},
	}
	cfg.AddCommand(cmd)
}
