package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
)

func registerPostGet(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "get <update-id>",
		Short: "Get a specific project update",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.ProjectPostGet(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			u := resp.ProjectUpdate
			result := mappers.FromProjectUpdateSummary(u.ProjectUpdateSummaryFields)
			result["project"] = map[string]any{
				"id":     u.Project.Id,
				"slugId": u.Project.SlugId,
				"name":   u.Project.Name,
			}

			output.PrintJSON(result)
		},
	})
}
