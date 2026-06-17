package customer

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	"github.com/shhac/lin/internal/output"
)

func Register(parent *cobra.Command) {
	customer := &cobra.Command{
		Use:   "customer",
		Short: "Customer and customer-request operations",
	}
	parent.AddCommand(customer)

	registerList(customer)
	registerSearch(customer)
	registerGet(customer)
	registerRequests(customer)
	registerStatuses(customer)
	registerTiers(customer)
	shared.RegisterUsage(customer, "customer", usageText)
	output.HandleUnknownCommand(customer, "To view a customer: lin customer get <id|slug>")
}

// referenceItem builds the output map shared by the customer statuses and
// tiers reference-list commands, which expose the same fields.
func referenceItem(id, name, displayName, color string, position float64, description *string) map[string]any {
	item := map[string]any{
		"id":          id,
		"name":        name,
		"displayName": displayName,
		"color":       color,
		"position":    position,
	}
	if description != nil {
		item["description"] = *description
	}
	return item
}
