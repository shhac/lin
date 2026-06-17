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
			items[i] = referenceItem(n.Id, n.Name, n.DisplayName, n.Color, n.Position, n.Description)
		}

		output.PrintPage(items, resp.CustomerTiers.PageInfo.HasNextPage, resp.CustomerTiers.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
