package customer

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
)

func registerSearch(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Search customers by name",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		filter := filters.BuildCustomerFilter(filters.CustomerFilterOpts{Search: args[0]})

		resp, err := linear.CustomerList(ctx, client, filter, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.Customers.Nodes))
		for i, n := range resp.Customers.Nodes {
			items[i] = mappers.MapCustomerSummary(n.CustomerSummaryFields)
		}

		output.PrintPage(items, resp.Customers.PageInfo.HasNextPage, resp.Customers.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
