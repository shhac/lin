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
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		Args:  cobra.NoArgs,
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
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

		resp, err := linear.ProjectList(ctx, client, filter, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.Projects.Nodes))
		for i, n := range resp.Projects.Nodes {
			items[i] = mappers.MapProjectSummary(mappers.FromProjectSummaryFields(n.ProjectSummaryFields))
		}

		output.PrintPage(items, resp.Projects.PageInfo.HasNextPage, resp.Projects.PageInfo.EndCursor)
	}

	cmd.Flags().StringVar(&team, "team", "", "Filter by team name or ID")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	parent.AddCommand(cmd)
}

