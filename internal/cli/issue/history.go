package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
)

func registerHistory(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "history <issue-id>",
		Short: "List activity history for an issue",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		resp, err := linear.IssueHistory(ctx, client, args[0], page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.Issue.History.Nodes))
		for i, h := range resp.Issue.History.Nodes {
			items[i] = mappers.MapHistoryEntry(h)
		}

		output.PrintPage(items, resp.Issue.History.PageInfo.HasNextPage, resp.Issue.History.PageInfo.EndCursor)
	}

	parent.AddCommand(cmd)
}
