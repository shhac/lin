package roadmap

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerProjects(roadmap *cobra.Command) {
	var limit string
	var cursor string

	cmd := &cobra.Command{
		Use:   "projects <id>",
		Short: "List projects linked to a roadmap",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			pageSize := output.ResolvePageSize(limit)

			resolved, err := resolvers.ResolveRoadmap(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			after := output.ResolveCursor(cursor)

			resp, err := linear.RoadmapProjects(context.Background(), client, resolved.ID, pageSize, after)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]map[string]any, len(resp.Roadmap.Projects.Nodes))
			for i, p := range resp.Roadmap.Projects.Nodes {
				f := p.ProjectSummaryFields
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

			pi := resp.Roadmap.Projects.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "50", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	roadmap.AddCommand(cmd)
}
