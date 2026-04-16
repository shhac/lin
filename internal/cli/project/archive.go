package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerDelete(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete (trash) a project",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectDelete(context.Background(), client, resolved.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"deleted": resp.ProjectDelete.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUnarchive(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "unarchive <id>",
		Short: "Restore a trashed or archived project",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUnarchive(context.Background(), client, resolved.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"unarchived": resp.ProjectUnarchive.Success})
		},
	}
	parent.AddCommand(cmd)
}
