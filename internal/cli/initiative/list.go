package initiative

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
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
	var limit string
	var cursor string
	var status string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List initiatives",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			client := linear.GetClient()
			pageSize := output.ResolvePageSize(limit)
			after := output.ResolveCursor(cursor)

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

			resp, err := linear.InitiativeList(context.Background(), client, filter, pageSize, after)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]map[string]any, len(resp.Initiatives.Nodes))
			for i, n := range resp.Initiatives.Nodes {
				var ownerName *string
				if n.Owner != nil {
					ownerName = &n.Owner.Name
				}
				item := map[string]any{
					"id":     n.Id,
					"slugId": n.SlugId,
					"url":    n.Url,
					"name":   n.Name,
					"status": n.Status,
					"owner":  ownerName,
				}
				if n.Health != nil {
					item["health"] = *n.Health
				}
				if n.TargetDate != nil {
					item["targetDate"] = *n.TargetDate
				}
				items[i] = item
			}

			pi := resp.Initiatives.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "Filter by status: planned|active|completed")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}
