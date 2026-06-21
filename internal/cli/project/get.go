package project

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/resolvers"
)

func registerGet(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>...",
		Short: "Project summary: status, progress, lead, dates, milestones",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return shared.GetEntities(args, func(client graphql.Client, id string) (any, error) {
				resolved, err := resolvers.ResolveProject(client, id)
				if err != nil {
					// ResolveProject returns plain fmt.Errorf; wrap so EntityGet
					// treats a missing project as an item-level @unresolved, not a
					// command-level failure.
					return nil, apierrors.Wrap(err, apierrors.FixableByAgent)
				}

				resp, err := linear.ProjectGet(context.Background(), client, resolved.ID)
				if err != nil {
					return nil, apierrors.ClassifyGraphQLError(err)
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

				return result, nil
			})
		},
	}

	parent.AddCommand(cmd)
}
