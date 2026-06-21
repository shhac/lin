package project

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
)

func registerPostGet(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "get <update-id>...",
		Short: "Get a specific project update",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return shared.GetEntities(args, func(client graphql.Client, id string) (any, error) {
				resp, err := linear.ProjectPostGet(context.Background(), client, id)
				if err != nil {
					return nil, apierrors.ClassifyGraphQLError(err)
				}

				u := resp.ProjectUpdate
				result := mappers.FromProjectUpdateSummary(u.ProjectUpdateSummaryFields)
				result["project"] = map[string]any{
					"id":     u.Project.Id,
					"slugId": u.Project.SlugId,
					"name":   u.Project.Name,
				}
				return result, nil
			})
		},
	})
}
