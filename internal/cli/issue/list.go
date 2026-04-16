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

func registerList(parent *cobra.Command) {
	var (
		project       string
		team          string
		assignee      string
		status        string
		priority      string
		label         string
		cycle         string
		updatedAfter  string
		updatedBefore string
		createdAfter  string
		createdBefore string
		limit         string
		cursor        string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			filter := filters.BuildIssueFilter(filters.IssueFilterOpts{
				Project:       project,
				Team:          team,
				Assignee:      assignee,
				Status:        status,
				Priority:      priority,
				Label:         label,
				Cycle:         cycle,
				UpdatedAfter:  updatedAfter,
				UpdatedBefore: updatedBefore,
				CreatedAfter:  createdAfter,
				CreatedBefore: createdBefore,
			})

			pageSize := output.ResolvePageSize(limit)
			afterPtr := output.ResolveCursor(cursor)

			resp, err := linear.IssueList(ctx, client, filter, pageSize, afterPtr)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.Issues.Nodes))
			for i, n := range resp.Issues.Nodes {
				items[i] = issueSummaryFromFields(n.IssueSummaryFields)
			}

			pi := resp.Issues.PageInfo
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
	cmd.Flags().StringVar(&label, "label", "", "Filter by label")
	cmd.Flags().StringVar(&cycle, "cycle", "", "Filter by cycle")
	cmd.Flags().StringVar(&updatedAfter, "updated-after", "", "Updated after date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&updatedBefore, "updated-before", "", "Updated before date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&createdAfter, "created-after", "", "Created after date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&createdBefore, "created-before", "", "Created before date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}

func issueSummaryFromFields(f linear.IssueSummaryFields) map[string]any {
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
	return mappers.MapIssueSummary(input)
}

