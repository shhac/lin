package customer

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
)

func registerGet(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id|slug>",
		Short: "Get customer details: tier, status, owner, domains, revenue, request count",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.CustomerGet(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(mappers.MapCustomerDetail(resp.Customer))
		},
	}

	parent.AddCommand(cmd)
}
