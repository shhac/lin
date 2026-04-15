package document

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
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
			var afterPtr *string
			if cursor != "" {
				afterPtr = &cursor
			}

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
				items[i] = mapSearchSummary(n.DocSearchSummaryFields)
			}

			pi := resp.SearchDocuments.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: derefStr(pi.EndCursor),
			})
		},
	}

	cmd.Flags().BoolVar(&includeComments, "include-comments", false, "Include comment text in search")
	cmd.Flags().BoolVar(&includeArchived, "include-archived", false, "Include archived documents")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}

func mapSearchSummary(f linear.DocSearchSummaryFields) map[string]any {
	m := map[string]any{
		"id":        f.Id,
		"slugId":    f.SlugId,
		"title":     f.Title,
		"url":       f.Url,
		"updatedAt": f.UpdatedAt,
	}
	if f.Project != nil {
		m["project"] = map[string]any{
			"id":   f.Project.Id,
			"name": f.Project.Name,
		}
	}
	if f.Creator != nil {
		m["creator"] = f.Creator.Name
	}
	return m
}
