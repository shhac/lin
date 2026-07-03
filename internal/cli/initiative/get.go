package initiative

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/output/pretty"
	"github.com/shhac/lin/internal/resolvers"
)

func registerGet(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>...",
		Short: "Initiative summary: name, description, status, health, owner",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			getOne := func(client graphql.Client, id string) (any, error) {
				resolved, err := resolvers.ResolveInitiative(client, id)
				if err != nil {
					return nil, err
				}

				resp, err := linear.InitiativeGet(context.Background(), client, resolved.ID)
				if err != nil {
					return nil, apierrors.ClassifyGraphQLError(err)
				}
				return mappers.MapInitiativeDetail(resp.Initiative), nil
			}
			if output.WantsPretty() {
				return shared.GetEntitiesPretty(args, getOne, func(item any, opts pretty.Options) string {
					return renderInitiativeCard(item.(map[string]any), opts)
				})
			}
			return shared.GetEntities(args, getOne)
		},
	}
	output.AllowPretty(cmd)
	parent.AddCommand(cmd)
}
