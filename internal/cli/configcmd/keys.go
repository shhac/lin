package configcmd

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output"
)

func registerListKeys(cfg *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "list-keys",
		Short: "List all available setting keys",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			keys := make([]map[string]string, 0, len(settingDefs))
			for key, def := range settingDefs {
				keys = append(keys, map[string]string{
					"key":         key,
					"description": def.description,
				})
			}
			output.PrintJSON(map[string]any{"keys": keys})
		},
	}
	cfg.AddCommand(cmd)
}
