package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
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
			var afterPtr *string
			if cursor != "" {
				afterPtr = &cursor
			}

			resp, err := linear.ProjectSearch(ctx, client, args[0], pageSize, afterPtr)
			if err != nil {
				output.PrintError(err.Error())
			}

			items := make([]any, len(resp.SearchProjects.Nodes))
			for i, n := range resp.SearchProjects.Nodes {
				items[i] = mapSearchSummary(n.ProjectSearchSummaryFields)
			}

			pi := resp.SearchProjects.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: derefStr(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}

func mapSearchSummary(f linear.ProjectSearchSummaryFields) map[string]any {
	m := map[string]any{
		"id":       f.Id,
		"slugId":   f.SlugId,
		"url":      f.Url,
		"name":     f.Name,
		"status":   f.State,
		"progress": f.Progress,
	}
	if f.Lead != nil {
		m["lead"] = f.Lead.Name
	}
	if f.StartDate != nil {
		m["startDate"] = *f.StartDate
	}
	if f.TargetDate != nil {
		m["targetDate"] = *f.TargetDate
	}
	return m
}
