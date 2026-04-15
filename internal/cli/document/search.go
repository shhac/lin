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

			var includeCommentsPtr *bool
			if includeComments {
				includeCommentsPtr = &includeComments
			}
			var includeArchivedPtr *bool
			if includeArchived {
				includeArchivedPtr = &includeArchived
			}

			resp, err := linear.DocumentSearch(ctx, client, args[0], pageSize, afterPtr, includeCommentsPtr, includeArchivedPtr)
			if err != nil {
				output.PrintError(err.Error())
			}

			items := make([]any, len(resp.SearchDocuments.Nodes))
			for i, n := range resp.SearchDocuments.Nodes {
				f := n.DocSearchSummaryFields
				input := mappers.DocSummaryInput{
					ID:        f.Id,
					SlugId:    f.SlugId,
					Title:     f.Title,
					URL:       f.Url,
					UpdatedAt: f.UpdatedAt,
				}
				if f.Creator != nil {
					input.CreatorID = f.Creator.Id
					input.CreatorName = f.Creator.Name
				}
				if f.Project != nil {
					input.ProjectID = f.Project.Id
					input.ProjectName = f.Project.Name
				}
				items[i] = mappers.MapDocSummary(input)
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

