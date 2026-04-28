package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerGet(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Project summary: status, progress, lead, dates, milestones",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectGet(ctx, client, resolved.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			p := resp.Project
			result := map[string]any{
				"id":          p.Id,
				"slugId":      p.SlugId,
				"url":         p.Url,
				"name":        p.Name,
				"description": p.Description,
				"content":     p.Content,
				"status":      p.State,
				"progress":    p.Progress,
				"startDate":   p.StartDate,
				"targetDate":  p.TargetDate,
			}

			if p.Lead != nil {
				result["lead"] = map[string]any{
					"id":   p.Lead.Id,
					"name": p.Lead.Name,
				}
			}

			labels := make([]map[string]any, len(p.Labels.Nodes))
			for i, l := range p.Labels.Nodes {
				labels[i] = map[string]any{"id": l.Id, "name": l.Name}
			}
			result["labels"] = labels

			milestones := make([]map[string]any, len(p.ProjectMilestones.Nodes))
			for i, m := range p.ProjectMilestones.Nodes {
				milestones[i] = map[string]any{
					"id":         m.Id,
					"name":       m.Name,
					"targetDate": m.TargetDate,
				}
			}
			result["milestones"] = milestones

			output.PrintJSON(result)
		},
	}

	parent.AddCommand(cmd)
}
