package document

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerNew(parent *cobra.Command) {
	var (
		project string
		content string
		icon    string
		color   string
	)

	cmd := &cobra.Command{
		Use:   "new <title>",
		Short: "Create a new document",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()
			title := args[0]

			input := linear.DocumentCreateInput{
				Title: title,
			}

			if project != "" {
				resolved, err := resolvers.ResolveProject(client, project)
				if err != nil {
					output.PrintError(err.Error())
				}
				input.ProjectId = &resolved.ID
			}
			if content != "" {
				input.Content = &content
			}
			if icon != "" {
				input.Icon = &icon
			}
			if color != "" {
				input.Color = &color
			}

			resp, err := linear.DocumentCreate(ctx, client, input)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			d := resp.DocumentCreate.Document
			output.PrintJSON(map[string]any{
				"id":      d.Id,
				"slugId":  d.SlugId,
				"title":   d.Title,
				"url":     d.Url,
				"created": resp.DocumentCreate.Success,
			})
		},
	}

	cmd.Flags().StringVar(&project, "project", "", "Project ID, slug, or name")
	cmd.Flags().StringVar(&content, "content", "", "Document content (markdown)")
	cmd.Flags().StringVar(&icon, "icon", "", "Document icon (emoji)")
	cmd.Flags().StringVar(&color, "color", "", "Icon color (hex)")
	parent.AddCommand(cmd)
}
