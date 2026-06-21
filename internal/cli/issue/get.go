package issue

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
)

func registerGet(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>...",
		Short: "Get issue details: title, description, status, assignee, labels, relationships",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return shared.GetEntities(args, func(client graphql.Client, id string) (any, error) {
				resp, err := linear.IssueGet(context.Background(), client, id)
				if err != nil {
					return nil, apierrors.ClassifyGraphQLError(err)
				}
				return mappers.MapIssueDetail(resp.Issue), nil
			})
		},
	}

	parent.AddCommand(cmd)
}
