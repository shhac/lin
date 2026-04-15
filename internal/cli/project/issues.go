package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerIssues(parent *cobra.Command) {
	var (
		status   string
		assignee string
		priority string
		limit    string
		cursor   string
	)

	cmd := &cobra.Command{
		Use:   "issues <id>",
		Short: "List issues within a project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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

			pageSize := output.ResolvePageSize(limit)
			var afterPtr *string
			if cursor != "" {
				afterPtr = &cursor
			}

			resp, err := linear.ProjectIssues(ctx, client, resolved.ID, filter, pageSize, afterPtr)
			if err != nil {
				output.PrintError(err.Error())
			}

			items := make([]any, len(resp.Project.Issues.Nodes))
			for i, n := range resp.Project.Issues.Nodes {
				items[i] = mapIssueSummary(n.IssueSummaryFields)
			}

			pi := resp.Project.Issues.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: derefStr(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Filter by assignee")
	cmd.Flags().StringVar(&priority, "priority", "", "Filter by priority")
	cmd.Flags().StringVar(&limit, "limit", "50", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}

func mapIssueSummary(f linear.IssueSummaryFields) map[string]any {
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
