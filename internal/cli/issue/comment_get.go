package issue

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
)

func registerCommentGet(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "get <comment-id>...",
		Short: "Get a specific comment",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return shared.GetEntities(args, func(client graphql.Client, id string) (any, error) {
				resp, err := linear.CommentGet(context.Background(), client, id)
				if err != nil {
					return nil, apierrors.ClassifyGraphQLError(err)
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
				return result, nil
			})
		},
	})
}
