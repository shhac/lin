package label

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerList(label *cobra.Command) {
	var teamFlag string
	var limit string
	var cursor string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			client := linear.GetClient()
			ctx := context.Background()
			pageSize := output.ResolvePageSize(limit)

			after := output.ResolveCursor(cursor)

			if teamFlag != "" {
				resolved, err := resolvers.ResolveTeam(client, teamFlag)
				if err != nil {
					output.PrintError(err.Error())
				}

				resp, err := linear.TeamLabels(ctx, client, resolved.ID, pageSize, after)
				if err != nil {
					output.HandleGraphQLError(err)
				}

				items := make([]map[string]any, len(resp.Team.Labels.Nodes))
				for i, l := range resp.Team.Labels.Nodes {
					items[i] = map[string]any{
						"id":    l.Id,
						"name":  l.Name,
						"color": l.Color,
					}
				}

				pi := resp.Team.Labels.PageInfo
				output.PrintPaginated(items, &output.Pagination{
					HasMore:    pi.HasNextPage,
					NextCursor: ptr.Deref(pi.EndCursor),
				})
				return
			}

			resp, err := linear.LabelList(ctx, client, pageSize, after)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]map[string]any, len(resp.IssueLabels.Nodes))
			for i, l := range resp.IssueLabels.Nodes {
				items[i] = map[string]any{
					"id":    l.Id,
					"name":  l.Name,
					"color": l.Color,
				}
			}

			pi := resp.IssueLabels.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&teamFlag, "team", "", "Filter by team")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	label.AddCommand(cmd)
}
