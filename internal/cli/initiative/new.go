package initiative

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerNew(parent *cobra.Command) {
	var (
		description string
		owner       string
		status      string
		targetDate  string
		content     string
		color       string
		icon        string
	)

	cmd := &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new initiative",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()
			name := args[0]

			input := linear.InitiativeCreateInput{
				Name: name,
			}

			if description != "" {
				input.Description = &description
			}
			if content != "" {
				input.Content = &content
			}
			if color != "" {
				input.Color = &color
			}
			if icon != "" {
				input.Icon = &icon
			}
			if targetDate != "" {
				input.TargetDate = &targetDate
			}

			if status != "" {
				normalized, err := validateInitiativeStatus(status)
				if err != nil {
					output.PrintError(err.Error())
				}
				s := linear.InitiativeStatus(normalized)
				input.Status = &s
			}

			if owner != "" {
				u, err := resolvers.ResolveUser(client, owner)
				if err != nil {
					output.PrintError(err.Error())
				}
				input.OwnerId = &u.ID
			}

			resp, err := linear.InitiativeCreate(ctx, client, input)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			i := resp.InitiativeCreate.Initiative
			output.PrintJSON(map[string]any{
				"id":      i.Id,
				"slugId":  i.SlugId,
				"url":     i.Url,
				"name":    i.Name,
				"created": resp.InitiativeCreate.Success,
			})
		},
	}

	cmd.Flags().StringVar(&description, "description", "", "Initiative description")
	cmd.Flags().StringVar(&owner, "owner", "", "Initiative owner: name, email, or user ID")
	cmd.Flags().StringVar(&status, "status", "", "Status: planned|active|completed")
	cmd.Flags().StringVar(&targetDate, "target-date", "", "Target date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&content, "content", "", "Initiative content body (markdown)")
	cmd.Flags().StringVar(&color, "color", "", "Initiative color")
	cmd.Flags().StringVar(&icon, "icon", "", "Initiative icon")
	parent.AddCommand(cmd)
}
