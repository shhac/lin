package initiative

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
		Short: "Search initiatives by name",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(_ *cobra.Command, args []string) {
		client := linear.GetClient()

		filter := &linear.InitiativeFilter{
			Name: filters.ContainsIgnoreCase(args[0]),
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

	parent.AddCommand(cmd)
}
