package project

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
					filter.State = &linear.StringComparator{EqIgnoreCase: ptr.To(status)}
				}
			}

			pageSize := output.ResolvePageSize(limit)
			afterPtr := output.ResolveCursor(cursor)

			resp, err := linear.ProjectList(ctx, client, filter, pageSize, afterPtr)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.Projects.Nodes))
			for i, n := range resp.Projects.Nodes {
				items[i] = mappers.MapProjectSummary(mappers.FromProjectSummaryFields(n.ProjectSummaryFields))
			}

			pi := resp.Projects.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&team, "team", "", "Filter by team name or ID")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}

