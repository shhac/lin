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
			items[i] = referenceItem(n.Id, n.Name, n.DisplayName, n.Color, n.Position, n.Description)
		}

		output.PrintPage(items, resp.CustomerStatuses.PageInfo.HasNextPage, resp.CustomerStatuses.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
