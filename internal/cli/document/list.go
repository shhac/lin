package document

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerList(parent *cobra.Command) {
	var (
		project         string
		creator         string
		includeArchived bool
		limit           string
		cursor          string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List documents",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			var filter *linear.DocumentFilter
			if project != "" || creator != "" {
				filter = &linear.DocumentFilter{}
				if project != "" {
					filter.Project = filters.BuildProjectFilter(project)
				}
				if creator != "" {
					filter.Creator = &linear.UserFilter{
						Or: []linear.UserFilter{
							{Id: &linear.IDComparator{Eq: strPtr(creator)}},
							{Name: &linear.StringComparator{EqIgnoreCase: strPtr(creator)}},
							{DisplayName: &linear.StringComparator{EqIgnoreCase: strPtr(creator)}},
							{Email: &linear.StringComparator{EqIgnoreCase: strPtr(creator)}},
						},
					}
				}
			}

			pageSize := output.ResolvePageSize(limit)
			var afterPtr *string
			if cursor != "" {
				afterPtr = &cursor
			}

			var includeArchivedPtr *bool
			if includeArchived {
				includeArchivedPtr = &includeArchived
			}

			resp, err := linear.DocumentList(ctx, client, filter, pageSize, afterPtr, includeArchivedPtr)
			if err != nil {
				output.PrintError(err.Error())
			}

			items := make([]any, len(resp.Documents.Nodes))
			for i, n := range resp.Documents.Nodes {
				items[i] = mapListSummary(n.DocSummaryFields)
			}

			pi := resp.Documents.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: derefStr(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&project, "project", "", "Filter by project ID, slug, or name")
	cmd.Flags().StringVar(&creator, "creator", "", "Filter by creator ID, name, or email")
	cmd.Flags().BoolVar(&includeArchived, "include-archived", false, "Include archived documents")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}

func mapListSummary(f linear.DocSummaryFields) map[string]any {
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
