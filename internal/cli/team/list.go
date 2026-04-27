package team

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerList(team *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all teams",
		Args:  cobra.NoArgs,
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(_ *cobra.Command, _ []string) {
		client := linear.GetClient()

		resp, err := linear.TeamList(context.Background(), client, nil, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]map[string]any, len(resp.Teams.Nodes))
		for i, t := range resp.Teams.Nodes {
			items[i] = map[string]any{
				"id":   t.Id,
				"name": t.Name,
				"key":  t.Key,
			}
		}

		output.PrintPage(items, resp.Teams.PageInfo.HasNextPage, resp.Teams.PageInfo.EndCursor)
	}

	team.AddCommand(cmd)
}
