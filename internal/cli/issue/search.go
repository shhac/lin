package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerSearch(parent *cobra.Command) {
	var (
		project  string
		team     string
		assignee string
		status   string
		priority string
		limit    string
		cursor   string
	)

	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Full-text search for issues",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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

			pageSize := output.ResolvePageSize(limit)
			var afterPtr *string
			if cursor != "" {
				afterPtr = &cursor
			}

			resp, err := linear.IssueSearch(ctx, client, text, pageSize, afterPtr, filter)
			if err != nil {
				output.PrintError(err.Error())
			}

			items := make([]any, len(resp.SearchIssues.Nodes))
			for i, n := range resp.SearchIssues.Nodes {
				items[i] = mapSearchSummary(n.IssueSearchSummaryFields)
			}

			pi := resp.SearchIssues.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: derefStr(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&project, "project", "", "Filter by project ID, slug, or name")
	cmd.Flags().StringVar(&team, "team", "", "Filter by team")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Filter by assignee")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&priority, "priority", "", "Filter by priority")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}

func mapSearchSummary(f linear.IssueSearchSummaryFields) map[string]any {
	m := map[string]any{
		"id":            f.Id,
		"identifier":    f.Identifier,
		"title":         f.Title,
		"branchName":    f.BranchName,
		"status":        f.State.Name,
		"statusType":    f.State.Type,
		"team":          f.Team.Key,
		"priority":      f.Priority,
		"priorityLabel": f.PriorityLabel,
	}
	if f.Assignee != nil {
		m["assignee"] = f.Assignee.Name
		m["assigneeId"] = f.Assignee.Id
	}
	return m
}
