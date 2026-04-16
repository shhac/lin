package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
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
			afterPtr := output.ResolveCursor(cursor)

			resp, err := linear.IssueSearch(ctx, client, text, pageSize, afterPtr, filter)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.SearchIssues.Nodes))
			for i, n := range resp.SearchIssues.Nodes {
				f := n.IssueSearchSummaryFields
				input := mappers.IssueSummaryInput{
					ID:            f.Id,
					Identifier:    f.Identifier,
					Title:         f.Title,
					BranchName:    f.BranchName,
					Priority:      f.Priority,
					PriorityLabel: f.PriorityLabel,
					StateName:     f.State.Name,
					StateType:     f.State.Type,
					TeamKey:       f.Team.Key,
				}
				if f.Assignee != nil {
					input.AssigneeID = f.Assignee.Id
					input.AssigneeName = f.Assignee.Name
				}
				items[i] = mappers.MapIssueSummary(input)
			}

			pi := resp.SearchIssues.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
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

