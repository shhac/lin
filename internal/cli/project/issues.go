package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerIssues(parent *cobra.Command) {
	var (
		status   string
		assignee string
		priority string
	)

	cmd := &cobra.Command{
		Use:   "issues <id>",
		Short: "List issues within a project",
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

		filter := filters.BuildIssueFilter(filters.IssueFilterOpts{
			Assignee: assignee,
			Status:   status,
			Priority: priority,
		})

		resp, err := linear.ProjectIssues(ctx, client, resolved.ID, filter, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.Project.Issues.Nodes))
		for i, n := range resp.Project.Issues.Nodes {
			items[i] = mappers.MapIssueSummary(mappers.FromIssueSummaryFields(n.IssueSummaryFields))
		}

		output.PrintPage(items, resp.Project.Issues.PageInfo.HasNextPage, resp.Project.Issues.PageInfo.EndCursor)
	}

	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Filter by assignee")
	cmd.Flags().StringVar(&priority, "priority", "", "Filter by priority")
	parent.AddCommand(cmd)
}
