package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerCommentGet(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "get <comment-id>",
		Short: "Get a specific comment",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.CommentGet(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			c := resp.Comment
			var uid, uname *string
			if c.User != nil {
				uid, uname = &c.User.Id, &c.User.Name
			}
			result := commentBase(c.Id, c.Body, c.CreatedAt, c.UpdatedAt, uid, uname)
			if c.Issue != nil {
				result["issue"] = map[string]any{"id": c.Issue.Id, "identifier": c.Issue.Identifier}
			}
			if c.Parent != nil {
				result["parent"] = map[string]any{"id": c.Parent.Id}
			}
			result["childCount"] = len(c.Children.Nodes)

			output.PrintJSON(result)
		},
	})
}
