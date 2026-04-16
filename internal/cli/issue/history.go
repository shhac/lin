package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
)

func registerHistory(parent *cobra.Command) {
	var (
		limit  string
		cursor string
	)

	cmd := &cobra.Command{
		Use:   "history <issue-id>",
		Short: "List activity history for an issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			pageSize := output.ResolvePageSize(limit)
			afterPtr := output.ResolveCursor(cursor)

			resp, err := linear.IssueHistory(ctx, client, args[0], pageSize, afterPtr)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.Issue.History.Nodes))
			for i, h := range resp.Issue.History.Nodes {
				items[i] = mappers.MapHistoryEntry(h)
			}

			pi := resp.Issue.History.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}
