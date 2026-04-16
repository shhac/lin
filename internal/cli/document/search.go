package document

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
)

func registerSearch(parent *cobra.Command) {
	var (
		includeComments bool
		includeArchived bool
		limit           string
		cursor          string
	)

	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Full-text search for documents",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			pageSize := output.ResolvePageSize(limit)
			afterPtr := output.ResolveCursor(cursor)

			resp, err := linear.DocumentSearch(ctx, client, args[0], pageSize, afterPtr, ptr.TrueOrNil(includeComments), ptr.TrueOrNil(includeArchived))
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.SearchDocuments.Nodes))
			for i, n := range resp.SearchDocuments.Nodes {
				items[i] = mappers.MapDocSummary(mappers.FromDocSearchSummaryFields(n.DocSearchSummaryFields))
			}

			pi := resp.SearchDocuments.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().BoolVar(&includeComments, "include-comments", false, "Include comment text in search")
	cmd.Flags().BoolVar(&includeArchived, "include-archived", false, "Include archived documents")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}

