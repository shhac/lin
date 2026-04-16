package project

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
		limit  string
		cursor string
	)

	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Full-text search for projects",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			pageSize := output.ResolvePageSize(limit)
			afterPtr := output.ResolveCursor(cursor)

			resp, err := linear.ProjectSearch(ctx, client, args[0], pageSize, afterPtr)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.SearchProjects.Nodes))
			for i, n := range resp.SearchProjects.Nodes {
				f := n.ProjectSearchSummaryFields
				input := mappers.ProjectSummaryInput{
					ID:         f.Id,
					SlugId:     f.SlugId,
					URL:        f.Url,
					Name:       f.Name,
					State:      f.State,
					Progress:   f.Progress,
					StartDate:  ptr.Deref(f.StartDate),
					TargetDate: ptr.Deref(f.TargetDate),
				}
				if f.Lead != nil {
					input.LeadName = f.Lead.Name
				}
				items[i] = mappers.MapProjectSummary(input)
			}

			pi := resp.SearchProjects.PageInfo
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

