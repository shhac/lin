package customer

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerRequests(parent *cobra.Command) {
	var (
		customer      string
		project       string
		status        string
		label         string
		team          string
		createdAfter  string
		createdBefore string
		important     bool
		unassigned    bool
		triage        bool
	)

	cmd := &cobra.Command{
		Use:   "requests",
		Short: "List customer requests across the workspace",
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		customerID := customer
		if customerID != "" {
			resolved, err := resolvers.ResolveCustomer(client, customer)
			if err != nil {
				output.PrintError(err.Error())
			}
			customerID = resolved.ID
		}

		filter := filters.BuildCustomerNeedFilter(filters.CustomerNeedFilterOpts{
			Customer:      customerID,
			Project:       project,
			Important:     important,
			Unassigned:    unassigned,
			Triage:        triage,
			Status:        status,
			Label:         label,
			Team:          team,
			CreatedAfter:  createdAfter,
			CreatedBefore: createdBefore,
		})

		resp, err := linear.CustomerNeeds(ctx, client, filter, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.CustomerNeeds.Nodes))
		for i, n := range resp.CustomerNeeds.Nodes {
			items[i] = mappers.MapCustomerNeedSummary(n.CustomerNeedSummaryFields)
		}

		output.PrintPage(items, resp.CustomerNeeds.PageInfo.HasNextPage, resp.CustomerNeeds.PageInfo.EndCursor)
	}

	cmd.Flags().StringVar(&customer, "customer", "", "Scope to a customer (UUID, slug, or name)")
	cmd.Flags().StringVar(&project, "project", "", "Scope to a project (UUID, slug, or name)")
	cmd.Flags().BoolVar(&important, "important", false, "Only important requests")
	cmd.Flags().BoolVar(&unassigned, "unassigned", false, "Only requests whose linked issue is unassigned")
	cmd.Flags().BoolVar(&triage, "triage", false, "Only requests whose linked issue is in triage")
	cmd.Flags().StringVar(&status, "status", "", "Filter by linked issue status name")
	cmd.Flags().StringVar(&label, "label", "", "Filter by linked issue label")
	cmd.Flags().StringVar(&team, "team", "", "Filter by linked issue team")
	cmd.Flags().StringVar(&createdAfter, "created-after", "", "Created after date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&createdBefore, "created-before", "", "Created before date (YYYY-MM-DD)")
	parent.AddCommand(cmd)
}
