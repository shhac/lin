package label

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerSearch(label *cobra.Command) {
	var (
		teamFlag string
		limit    string
		cursor   string
	)

	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Search labels by name (case- and accent-insensitive substring)",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()
			pageSize := output.ResolvePageSize(limit)
			after := output.ResolveCursor(cursor)

			var teamID string
			if teamFlag != "" {
				resolved, err := resolvers.ResolveTeam(client, teamFlag)
				if err != nil {
					output.PrintError(err.Error())
				}
				teamID = resolved.ID
			}

			filter := filters.BuildIssueLabelFilter(filters.LabelFilterOpts{Search: args[0]}, teamID)

			resp, err := linear.LabelList(ctx, client, pageSize, after, filter)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]map[string]any, len(resp.IssueLabels.Nodes))
			for i, n := range resp.IssueLabels.Nodes {
				items[i] = mapLabel(n.LabelFields)
			}

			pi := resp.IssueLabels.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&teamFlag, "team", "", "Restrict search to a single team (key, name, or UUID)")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	label.AddCommand(cmd)
}
