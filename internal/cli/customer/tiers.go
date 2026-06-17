package customer

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerTiers(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "tiers",
		Short: "List workspace customer tiers (segments)",
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		resp, err := linear.CustomerTiers(ctx, client, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.CustomerTiers.Nodes))
		for i, n := range resp.CustomerTiers.Nodes {
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

		output.PrintPage(items, resp.CustomerTiers.PageInfo.HasNextPage, resp.CustomerTiers.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
