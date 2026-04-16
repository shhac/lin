package document

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
		Short: "Get document details (includes full content)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveDocument(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.DocumentGet(ctx, client, resolved.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			d := resp.Document
			result := map[string]any{
				"id":        d.Id,
				"slugId":    d.SlugId,
				"title":     d.Title,
				"content":   d.Content,
				"url":       d.Url,
				"icon":      d.Icon,
				"color":     d.Color,
				"createdAt": d.CreatedAt,
				"updatedAt": d.UpdatedAt,
			}

			if d.Project != nil {
				result["project"] = map[string]any{
					"id":     d.Project.Id,
					"name":   d.Project.Name,
					"slugId": d.Project.SlugId,
				}
			}
			if d.Creator != nil {
				result["creator"] = map[string]any{
					"id":   d.Creator.Id,
					"name": d.Creator.Name,
				}
			}
			if d.UpdatedBy != nil {
				result["updatedBy"] = map[string]any{
					"id":   d.UpdatedBy.Id,
					"name": d.UpdatedBy.Name,
				}
			}

			output.PrintJSON(result)
		},
	}

	parent.AddCommand(cmd)
}
