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

			resp, err := linear.DocumentList(ctx, client, filter, pageSize, afterPtr, ptr.TrueOrNil(includeArchived))
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.Documents.Nodes))
			for i, n := range resp.Documents.Nodes {
				items[i] = mappers.MapDocSummary(mappers.FromDocSummaryFields(n.DocSummaryFields))
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

