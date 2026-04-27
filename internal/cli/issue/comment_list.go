package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerCommentList(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "list <issue-id>",
		Short: "List comments on an issue",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		resp, err := linear.IssueComments(ctx, client, args[0], page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.Issue.Comments.Nodes))
		for i, c := range resp.Issue.Comments.Nodes {
			var uid, uname *string
			if c.User != nil {
				uid, uname = &c.User.Id, &c.User.Name
			}
			m := commentBase(c.Id, c.Body, c.CreatedAt, c.UpdatedAt, uid, uname)
			if c.Parent != nil {
				m["parent"] = map[string]any{"id": c.Parent.Id}
			}
			m["childCount"] = len(c.Children.Nodes)
			items[i] = m
		}

		output.PrintPage(items, resp.Issue.Comments.PageInfo.HasNextPage, resp.Issue.Comments.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
