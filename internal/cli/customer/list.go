package customer

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
)

func registerList(parent *cobra.Command) {
	var (
		tier    string
		status  string
		owner   string
		domain  string
		revenue string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List customers (filter by tier, status, owner, domain, revenue)",
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		filter := filters.BuildCustomerFilter(filters.CustomerFilterOpts{
			Tier:    tier,
			Status:  status,
			Owner:   owner,
			Domain:  domain,
			Revenue: revenue,
		})

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

	cmd.Flags().StringVar(&tier, "tier", "", "Filter by tier display name")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status name")
	cmd.Flags().StringVar(&owner, "owner", "", "Filter by owner (me, name, email, or UUID)")
	cmd.Flags().StringVar(&domain, "domain", "", "Filter by email domain")
	cmd.Flags().StringVar(&revenue, "revenue", "", "Minimum annual revenue")
	parent.AddCommand(cmd)
}
