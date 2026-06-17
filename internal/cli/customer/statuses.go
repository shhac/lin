package customer

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerStatuses(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "statuses",
		Short: "List workspace customer lifecycle statuses",
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		resp, err := linear.CustomerStatuses(ctx, client, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.CustomerStatuses.Nodes))
		for i, n := range resp.CustomerStatuses.Nodes {
			item := map[string]any{
				"id":          n.Id,
				"name":        n.Name,
				"displayName": n.DisplayName,
				"color":       n.Color,
				"position":    n.Position,
			}
			if n.Description != nil {
				item["description"] = *n.Description
			}
			items[i] = item
		}

		output.PrintPage(items, resp.CustomerStatuses.PageInfo.HasNextPage, resp.CustomerStatuses.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
