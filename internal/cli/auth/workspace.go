package auth

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/output"
)

func registerWorkspace(auth *cobra.Command) {
	workspace := &cobra.Command{
		Use:   "workspace",
		Short: "Manage workspace profiles",
	}
	output.HandleUnknownCommand(workspace, "Run 'lin auth workspace list' for available workspaces")

	registerWorkspaceList(workspace)
	registerWorkspaceSwitch(workspace)
	registerWorkspaceRemove(workspace)

	auth.AddCommand(workspace)
}

func registerWorkspaceList(workspace *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all stored workspaces",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			workspaces := config.GetWorkspaces()
			defaultWs := config.GetDefaultWorkspace()
			items := make([]map[string]any, 0, len(workspaces))
			for alias, ws := range workspaces {
				items = append(items, map[string]any{
					"alias":   alias,
					"name":    ws.Name,
					"urlKey":  ws.URLKey,
					"default": alias == defaultWs,
				})
			}
			output.PrintJSON(map[string]any{"items": items})
		},
	}
	workspace.AddCommand(cmd)
}

func registerWorkspaceSwitch(workspace *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "switch <alias>",
		Short: "Set default workspace",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			alias := args[0]
			if err := config.SetDefaultWorkspace(alias); err != nil {
				output.PrintError(err.Error())
				return
			}
			output.PrintJSON(map[string]any{"ok": true, "default_workspace": alias})
		},
	}
	workspace.AddCommand(cmd)
}

func registerWorkspaceRemove(workspace *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove a stored workspace",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			alias := args[0]
			wasDefault := config.GetDefaultWorkspace() == alias
			if err := config.RemoveWorkspace(alias); err != nil {
				output.PrintError(err.Error())
				return
			}
			newDefault := config.GetDefaultWorkspace()
			result := map[string]any{
				"ok":                true,
				"removed":           alias,
				"default_workspace": nilIfEmpty(newDefault),
			}
			if wasDefault {
				result["warning"] = "Removed the default workspace"
			}
			output.PrintJSON(result)
		},
	}
	workspace.AddCommand(cmd)
}
