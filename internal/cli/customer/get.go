package customer

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
)

func registerGet(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id|slug>...",
		Short: "Get customer details: tier, status, owner, domains, revenue, request count",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			getOne := func(client graphql.Client, id string) (any, error) {
				resp, err := linear.CustomerGet(context.Background(), client, id)
				if err != nil {
					return nil, apierrors.ClassifyGraphQLError(err)
				}
				return mappers.MapCustomerDetail(resp.Customer), nil
			}
			if output.WantsPretty() {
				return shared.GetEntitiesPretty(args, getOne, func(item any, opts pretty.Options) string {
					return renderCustomerCard(item.(map[string]any), opts)
				})
			}
			return shared.GetEntities(args, getOne)
		},
	}
	output.AllowPretty(cmd)

	parent.AddCommand(cmd)
}
