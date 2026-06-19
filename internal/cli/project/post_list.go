package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerPostList(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "list <project>",
		Short: "List project updates (newest first)",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		resolved, err := resolvers.ResolveProject(client, args[0])
		if err != nil {
			output.PrintError(err.Error())
		}

		resp, err := linear.ProjectPostList(ctx, client, resolved.ID, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		updates := resp.Project.ProjectUpdates
		items := make([]any, len(updates.Nodes))
		for i, n := range updates.Nodes {
			items[i] = mappers.FromProjectUpdateSummary(n.ProjectUpdateSummaryFields)
		}

		output.PrintPage(items, updates.PageInfo.HasNextPage, updates.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
