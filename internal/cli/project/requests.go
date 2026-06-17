package project

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
	var important bool

	cmd := &cobra.Command{
		Use:   "requests <id>",
		Short: "List customer requests linked to a project",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		resolved, err := resolvers.ResolveProject(client, args[0])
		if err != nil {
			output.PrintError(err.Error())
		}

		filter := filters.BuildCustomerNeedFilter(filters.CustomerNeedFilterOpts{Important: important})

		resp, err := linear.ProjectNeeds(ctx, client, resolved.ID, filter, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		conn := resp.Project.Needs
		items := make([]any, len(conn.Nodes))
		for i, n := range conn.Nodes {
			items[i] = mappers.MapCustomerNeedSummary(n.CustomerNeedSummaryFields)
		}

		output.PrintPage(items, conn.PageInfo.HasNextPage, conn.PageInfo.EndCursor)
	}

	cmd.Flags().BoolVar(&important, "important", false, "Only important requests")
	parent.AddCommand(cmd)
}
