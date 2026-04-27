package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdateLead(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "lead <id> <user>",
		Short: "Update project lead",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			user, err := resolvers.ResolveUser(client, args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				LeadId: &user.ID,
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}
