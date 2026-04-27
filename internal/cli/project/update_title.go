package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdateTitle(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "title <id> <new-title>",
		Short: "Update project title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				Name: &args[1],
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{
				"id":      resolved.ID,
				"name":    args[1],
				"updated": resp.ProjectUpdate.Success,
			})
		},
	}
	parent.AddCommand(cmd)
}
