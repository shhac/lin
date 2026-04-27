package configcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/output"
)

func registerSet(cfg *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Update a setting",
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			key, value := args[0], args[1]
			def, ok := settingDefs[key]
			if !ok {
				output.PrintError(fmt.Sprintf("Unknown setting: %s. Valid keys: %s", key, validKeysStr()))
				return
			}

			parsed, err := def.parse(value)
			if err != nil {
				output.PrintError(err.Error())
				return
			}

			intVal := parsed.(int)
			settings := config.GetSettings()
			def.set(settings, intVal)

			if err := config.UpdateSettings(settings); err != nil {
				output.PrintError(err.Error())
				return
			}
			output.PrintJSON(map[string]any{key: parsed})
		},
	}
	cfg.AddCommand(cmd)
}
