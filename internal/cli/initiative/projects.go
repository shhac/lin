package initiative

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerProjects(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "projects <id>",
		Short: "List projects linked to an initiative",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(_ *cobra.Command, args []string) {
		client := linear.GetClient()

		resolved, err := resolvers.ResolveInitiative(client, args[0])
		if err != nil {
			output.PrintError(err.Error())
		}

		resp, err := linear.InitiativeProjects(context.Background(), client, resolved.ID, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]map[string]any, len(resp.Initiative.Projects.Nodes))
		for i, p := range resp.Initiative.Projects.Nodes {
			items[i] = mappers.MapProjectSummary(mappers.FromProjectSummaryFields(p.ProjectSummaryFields))
		}

		output.PrintPage(items, resp.Initiative.Projects.PageInfo.HasNextPage, resp.Initiative.Projects.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
