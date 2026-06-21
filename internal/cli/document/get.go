package document

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/resolvers"
)

func registerGet(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>...",
		Short: "Get document details (includes full content)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return shared.GetEntities(args, func(client graphql.Client, id string) (any, error) {
				resolved, err := resolvers.ResolveDocument(client, id)
				if err != nil {
					// ResolveDocument returns plain fmt.Errorf; wrap so EntityGet
					// treats a missing document as an item-level @unresolved, not a
					// command-level failure.
					return nil, apierrors.Wrap(err, apierrors.FixableByAgent)
				}

				resp, err := linear.DocumentGet(context.Background(), client, resolved.ID)
				if err != nil {
					return nil, apierrors.ClassifyGraphQLError(err)
				}

				d := resp.Document
				result := map[string]any{
					"id":        d.Id,
					"slugId":    d.SlugId,
					"title":     d.Title,
					"content":   d.Content,
					"url":       d.Url,
					"icon":      d.Icon,
					"color":     d.Color,
					"createdAt": d.CreatedAt,
					"updatedAt": d.UpdatedAt,
				}

				if d.Project != nil {
					result["project"] = map[string]any{
						"id":     d.Project.Id,
						"name":   d.Project.Name,
						"slugId": d.Project.SlugId,
					}
				}
				if d.Creator != nil {
					result["creator"] = map[string]any{
						"id":   d.Creator.Id,
						"name": d.Creator.Name,
					}
				}
				if d.UpdatedBy != nil {
					result["updatedBy"] = map[string]any{
						"id":   d.UpdatedBy.Id,
						"name": d.UpdatedBy.Name,
					}
				}
				return result, nil
			})
		},
	}

	parent.AddCommand(cmd)
}
