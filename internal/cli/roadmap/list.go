package roadmap

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
)

func registerList(roadmap *cobra.Command) {
	var limit string
	var cursor string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List roadmaps",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			client := linear.GetClient()
			pageSize := output.ResolvePageSize(limit)

			after := output.ResolveCursor(cursor)

			resp, err := linear.RoadmapList(context.Background(), client, pageSize, after)
			if err != nil {
				output.PrintError(err.Error())
			}

			items := make([]map[string]any, len(resp.Roadmaps.Nodes))
			for i, r := range resp.Roadmaps.Nodes {
				var ownerName *string
				if r.Owner != nil {
					ownerName = &r.Owner.Name
				}
				items[i] = map[string]any{
					"id":          r.Id,
					"slugId":      r.SlugId,
					"url":         r.Url,
					"name":        r.Name,
					"description": r.Description,
					"owner":       ownerName,
				}
			}

			pi := resp.Roadmaps.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	roadmap.AddCommand(cmd)
}

