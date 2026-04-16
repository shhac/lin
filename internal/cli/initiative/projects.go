package initiative

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerProjects(parent *cobra.Command) {
	var limit string
	var cursor string

	cmd := &cobra.Command{
		Use:   "projects <id>",
		Short: "List projects linked to an initiative",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			pageSize := output.ResolvePageSize(limit)

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			after := output.ResolveCursor(cursor)

			resp, err := linear.InitiativeProjects(context.Background(), client, resolved.ID, pageSize, after)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]map[string]any, len(resp.Initiative.Projects.Nodes))
			for i, p := range resp.Initiative.Projects.Nodes {
				items[i] = mappers.MapProjectSummary(mappers.FromProjectSummaryFields(p.ProjectSummaryFields))
			}

			pi := resp.Initiative.Projects.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "50", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}
