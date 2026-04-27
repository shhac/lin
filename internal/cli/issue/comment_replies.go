package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerCommentReplies(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "replies <comment-id>",
		Short: "List replies to a comment",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		resp, err := linear.CommentReplies(ctx, client, args[0], page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.Comment.Children.Nodes))
		for i, c := range resp.Comment.Children.Nodes {
			var uid, uname *string
			if c.User != nil {
				uid, uname = &c.User.Id, &c.User.Name
			}
			items[i] = commentBase(c.Id, c.Body, c.CreatedAt, c.UpdatedAt, uid, uname)
		}

		output.PrintPage(items, resp.Comment.Children.PageInfo.HasNextPage, resp.Comment.Children.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
