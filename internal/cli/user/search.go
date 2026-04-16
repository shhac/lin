package user

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
)

func registerSearch(user *cobra.Command) {
	var limit string
	var cursor string

	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Search users by name, email, or display name",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			text := args[0]
			client := linear.GetClient()
			pageSize := output.ResolvePageSize(limit)

			after := output.ResolveCursor(cursor)

			filter := &linear.UserFilter{
				Or: []linear.UserFilter{
					{Name: containsIgnoreCase(text)},
					{DisplayName: containsIgnoreCase(text)},
					{Email: containsIgnoreCase(text)},
				},
			}

			resp, err := linear.UserList(context.Background(), client, filter, pageSize, after)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]map[string]any, len(resp.Users.Nodes))
			for i, u := range resp.Users.Nodes {
				items[i] = mapUserSummary(u)
			}

			pi := resp.Users.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	user.AddCommand(cmd)
}

func containsIgnoreCase(s string) *linear.StringComparator {
	return &linear.StringComparator{ContainsIgnoreCase: ptr.To(s)}
}
