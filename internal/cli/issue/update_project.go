package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdateProject(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "project <id> <project>",
		Short: "Move issue to project",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			proj, err := resolvers.ResolveProject(client, args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{ProjectId: &proj.ID})
			if err != nil {
				output.HandleGraphQLError(err)
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}
