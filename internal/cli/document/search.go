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
	)

	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Full-text search for documents",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		resp, err := linear.DocumentSearch(ctx, client, args[0], page.Size(), page.Cursor(), ptr.TrueOrNil(includeComments), ptr.TrueOrNil(includeArchived))
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]any, len(resp.SearchDocuments.Nodes))
		for i, n := range resp.SearchDocuments.Nodes {
			items[i] = mappers.MapDocSummary(mappers.FromDocSearchSummaryFields(n.DocSearchSummaryFields))
		}

		output.PrintPage(items, resp.SearchDocuments.PageInfo.HasNextPage, resp.SearchDocuments.PageInfo.EndCursor)
	}

	cmd.Flags().BoolVar(&includeComments, "include-comments", false, "Include comment text in search")
	cmd.Flags().BoolVar(&includeArchived, "include-archived", false, "Include archived documents")
	parent.AddCommand(cmd)
}
