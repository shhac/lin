package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/projectstatuses"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdateStatus(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "status <id> <new-status>",
		Short: "Update project status",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			normalized, err := projectstatuses.Validate(args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				State: &normalized,
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}
