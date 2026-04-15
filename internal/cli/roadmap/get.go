package roadmap

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerGet(roadmap *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Roadmap summary: name, description, owner",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()

			resolved, err := resolvers.ResolveRoadmap(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.RoadmapGet(context.Background(), client, resolved.ID)
			if err != nil {
				output.PrintError(err.Error())
			}

			r := resp.Roadmap

			var owner any
			if r.Owner != nil {
				owner = map[string]any{
					"id":   r.Owner.Id,
					"name": r.Owner.Name,
				}
			}

			output.PrintJSON(map[string]any{
				"id":          r.Id,
				"slugId":      r.SlugId,
				"url":         r.Url,
				"name":        r.Name,
				"description": r.Description,
				"owner":       owner,
				"creator": map[string]any{
					"id":   r.Creator.Id,
					"name": r.Creator.Name,
				},
				"createdAt": r.CreatedAt,
			})
		},
	}
	roadmap.AddCommand(cmd)
}
