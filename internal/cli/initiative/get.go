package initiative

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
		Short: "Initiative summary: name, description, status, health, owner",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.InitiativeGet(context.Background(), client, resolved.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			i := resp.Initiative

			var owner any
			if i.Owner != nil {
				owner = map[string]any{
					"id":   i.Owner.Id,
					"name": i.Owner.Name,
				}
			}

			result := map[string]any{
				"id":        i.Id,
				"slugId":    i.SlugId,
				"url":       i.Url,
				"name":      i.Name,
				"status":    i.Status,
				"owner":     owner,
				"creator": map[string]any{
					"id":   i.Creator.Id,
					"name": i.Creator.Name,
				},
				"createdAt": i.CreatedAt,
				"updatedAt": i.UpdatedAt,
			}
			if i.Description != nil {
				result["description"] = *i.Description
			}
			if i.Content != nil {
				result["content"] = *i.Content
			}
			if i.Health != nil {
				result["health"] = *i.Health
			}
			if i.TargetDate != nil {
				result["targetDate"] = *i.TargetDate
			}
			if i.StartedAt != nil {
				result["startedAt"] = *i.StartedAt
			}
			if i.CompletedAt != nil {
				result["completedAt"] = *i.CompletedAt
			}

			output.PrintJSON(result)
		},
	}
	parent.AddCommand(cmd)
}
