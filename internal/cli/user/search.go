package user

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
)

func registerSearch(user *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Search users by name, email, or display name",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(_ *cobra.Command, args []string) {
		text := args[0]
		client := linear.GetClient()

		filter := &linear.UserFilter{
			Or: []linear.UserFilter{
				{Name: containsIgnoreCase(text)},
				{DisplayName: containsIgnoreCase(text)},
				{Email: containsIgnoreCase(text)},
			},
		}

		resp, err := linear.UserList(context.Background(), client, filter, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]map[string]any, len(resp.Users.Nodes))
		for i, u := range resp.Users.Nodes {
			items[i] = mapUserSummary(u)
		}

		output.PrintPage(items, resp.Users.PageInfo.HasNextPage, resp.Users.PageInfo.EndCursor)
	}

	user.AddCommand(cmd)
}

func containsIgnoreCase(s string) *linear.StringComparator {
	return &linear.StringComparator{ContainsIgnoreCase: ptr.To(s)}
}
