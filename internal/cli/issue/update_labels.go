package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdateLabels(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "labels <id> <labels>",
		Short: "Set issue labels",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			teamID := resolveIssueTeamID(ctx, client, args[0])

			labelIds, err := resolvers.ResolveIssueLabels(client, args[1], teamID)
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{LabelIds: labelIds})
			if err != nil {
				output.HandleGraphQLError(err)
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}
