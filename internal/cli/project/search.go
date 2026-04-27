package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
)

func registerSearch(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Full-text search for projects",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		resp, err := linear.ProjectSearch(ctx, client, args[0], page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.SearchProjects.Nodes))
		for i, n := range resp.SearchProjects.Nodes {
			items[i] = mappers.MapProjectSummary(mappers.FromProjectSearchSummaryFields(n.ProjectSearchSummaryFields))
		}

		output.PrintPage(items, resp.SearchProjects.PageInfo.HasNextPage, resp.SearchProjects.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
