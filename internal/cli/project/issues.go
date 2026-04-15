package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
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
			afterPtr := output.ResolveCursor(cursor)

			resp, err := linear.ProjectIssues(ctx, client, resolved.ID, filter, pageSize, afterPtr)
			if err != nil {
				output.PrintError(err.Error())
			}

			items := make([]any, len(resp.Project.Issues.Nodes))
			for i, n := range resp.Project.Issues.Nodes {
				f := n.IssueSummaryFields
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

			pi := resp.Project.Issues.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
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

