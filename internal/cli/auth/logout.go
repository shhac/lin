package auth

import (
	"sort"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/output"
)

func registerLogout(auth *cobra.Command) {
	var all bool

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Clear stored credentials",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			if all {
				if err := config.ClearAll(); err != nil {
					output.PrintError(err.Error())
					return
				}
				output.PrintJSON(map[string]any{"ok": true, "cleared": "all"})
				return
			}

			defaultWs := config.GetDefaultWorkspace()
			if defaultWs != "" {
				_ = config.RemoveWorkspace(defaultWs)
			}
			_ = config.ClearApiKey()

			newDefault := config.GetDefaultWorkspace()
			remaining := sortedKeys(config.GetWorkspaces())

			result := map[string]any{
				"ok":                   true,
				"removed":             nilIfEmpty(defaultWs),
				"remaining_workspaces": remaining,
				"default_workspace":   nilIfEmpty(newDefault),
			}
			output.PrintJSON(result)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Remove all workspaces (default: only active workspace)")
	auth.AddCommand(cmd)
}

func sortedKeys(m map[string]config.Workspace) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
