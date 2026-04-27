package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
)

func registerSearch(parent *cobra.Command) {
	var (
		project  string
		team     string
		assignee string
		status   string
		priority string
	)

	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Full-text search for issues",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()
		text := args[0]

		filter := filters.BuildIssueFilter(filters.IssueFilterOpts{
			Project:  project,
			Team:     team,
			Assignee: assignee,
			Status:   status,
			Priority: priority,
		})

		resp, err := linear.IssueSearch(ctx, client, text, page.Size(), page.Cursor(), filter)
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.SearchIssues.Nodes))
		for i, n := range resp.SearchIssues.Nodes {
			items[i] = mappers.MapIssueSummary(mappers.FromIssueSearchSummaryFields(n.IssueSearchSummaryFields))
		}

		output.PrintPage(items, resp.SearchIssues.PageInfo.HasNextPage, resp.SearchIssues.PageInfo.EndCursor)
	}

	cmd.Flags().StringVar(&project, "project", "", "Filter by project ID, slug, or name")
	cmd.Flags().StringVar(&team, "team", "", "Filter by team")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Filter by assignee")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&priority, "priority", "", "Filter by priority")
	parent.AddCommand(cmd)
}
