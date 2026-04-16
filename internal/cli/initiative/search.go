package initiative

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
)

func registerSearch(parent *cobra.Command) {
	var (
		limit  string
		cursor string
	)

	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Search initiatives by name",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			pageSize := output.ResolvePageSize(limit)
			after := output.ResolveCursor(cursor)

			filter := &linear.InitiativeFilter{
				Name: &linear.StringComparator{ContainsIgnoreCase: &args[0]},
			}

			resp, err := linear.InitiativeList(context.Background(), client, filter, pageSize, after)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]map[string]any, len(resp.Initiatives.Nodes))
			for i, n := range resp.Initiatives.Nodes {
				items[i] = mappers.MapInitiativeSummary(mappers.FromInitiativeListFields(n))
			}

			pi := resp.Initiatives.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}
