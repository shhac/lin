package document

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
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
							{Id: &linear.IDComparator{Eq: ptr.To(creator)}},
							{Name: &linear.StringComparator{EqIgnoreCase: ptr.To(creator)}},
							{DisplayName: &linear.StringComparator{EqIgnoreCase: ptr.To(creator)}},
							{Email: &linear.StringComparator{EqIgnoreCase: ptr.To(creator)}},
						},
					}
				}
			}

			pageSize := output.ResolvePageSize(limit)
			afterPtr := output.ResolveCursor(cursor)

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
				items[i] = docSummaryFromFields(n.DocSummaryFields)
			}

			pi := resp.Documents.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
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

func docSummaryFromFields(f linear.DocSummaryFields) map[string]any {
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
	return mappers.MapDocSummary(input)
}
