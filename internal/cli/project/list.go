package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerList(parent *cobra.Command) {
	var (
		team   string
		status string
		limit  string
		cursor string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			var filter *linear.ProjectFilter
			if team != "" || status != "" {
				filter = &linear.ProjectFilter{}
				if team != "" {
					filter.AccessibleTeams = &linear.TeamCollectionFilter{
						Some: filters.BuildTeamFilter(team),
					}
				}
				if status != "" {
					filter.State = &linear.StringComparator{EqIgnoreCase: strPtr(status)}
				}
			}

			pageSize := output.ResolvePageSize(limit)
			var afterPtr *string
			if cursor != "" {
				afterPtr = &cursor
			}

			resp, err := linear.ProjectList(ctx, client, filter, pageSize, afterPtr)
			if err != nil {
				output.PrintError(err.Error())
			}

			items := make([]any, len(resp.Projects.Nodes))
			for i, n := range resp.Projects.Nodes {
				items[i] = mapListSummary(n.ProjectSummaryFields)
			}

			pi := resp.Projects.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: derefStr(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&team, "team", "", "Filter by team name or ID")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}

func mapListSummary(f linear.ProjectSummaryFields) map[string]any {
	m := map[string]any{
		"id":       f.Id,
		"slugId":   f.SlugId,
		"url":      f.Url,
		"name":     f.Name,
		"status":   f.State,
		"progress": f.Progress,
	}
	if f.Lead != nil {
		m["lead"] = f.Lead.Name
	}
	if f.StartDate != nil {
		m["startDate"] = *f.StartDate
	}
	if f.TargetDate != nil {
		m["targetDate"] = *f.TargetDate
	}
	return m
}
