package initiative

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
)

var validStatuses = []string{"planned", "active", "completed"}

const statusValues = "planned | active | completed"

func validateInitiativeStatus(input string) (string, error) {
	lower := strings.ToLower(input)
	for _, s := range validStatuses {
		if s == lower {
			// Linear expects capitalized status values
			return strings.ToUpper(s[:1]) + s[1:], nil
		}
	}
	return "", fmt.Errorf("unknown initiative status: %q, valid values: %s", input, statusValues)
}

func registerList(parent *cobra.Command) {
	var status string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List initiatives",
		Args:  cobra.NoArgs,
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(_ *cobra.Command, _ []string) {
		client := linear.GetClient()

		var filter *linear.InitiativeFilter
		if status != "" {
			normalized, err := validateInitiativeStatus(status)
			if err != nil {
				output.PrintError(err.Error())
			}
			filter = &linear.InitiativeFilter{
				Status: &linear.StringComparator{EqIgnoreCase: &normalized},
			}
		}

		resp, err := linear.InitiativeList(context.Background(), client, filter, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]map[string]any, len(resp.Initiatives.Nodes))
		for i, n := range resp.Initiatives.Nodes {
			items[i] = mappers.MapInitiativeSummary(mappers.FromInitiativeListFields(n))
		}

		output.PrintPage(items, resp.Initiatives.PageInfo.HasNextPage, resp.Initiatives.PageInfo.EndCursor)
	}

	cmd.Flags().StringVar(&status, "status", "", "Filter by status: planned|active|completed")
	parent.AddCommand(cmd)
}
