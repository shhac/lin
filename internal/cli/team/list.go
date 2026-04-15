package team

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
)

func registerList(team *cobra.Command) {
	var limit string
	var cursor string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all teams",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			client := linear.GetClient()
			pageSize := output.ResolvePageSize(limit)

			after := output.ResolveCursor(cursor)

			resp, err := linear.TeamList(context.Background(), client, nil, pageSize, after)
			if err != nil {
				output.PrintError(err.Error())
			}

			items := make([]map[string]any, len(resp.Teams.Nodes))
			for i, t := range resp.Teams.Nodes {
				items[i] = map[string]any{
					"id":   t.Id,
					"name": t.Name,
					"key":  t.Key,
				}
			}

			pi := resp.Teams.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	team.AddCommand(cmd)
}

