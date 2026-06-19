package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/projecthealth"
	"github.com/shhac/lin/internal/resolvers"
)

func registerPostNew(parent *cobra.Command) {
	var health string

	cmd := &cobra.Command{
		Use:   "new <project> <body>",
		Short: "Post a project update",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			body := args[1]
			input := linear.ProjectUpdateCreateInput{
				ProjectId: resolved.ID,
				Body:      &body,
			}
			if health != "" {
				h, err := projecthealth.Validate(health)
				if err != nil {
					output.PrintError(err.Error())
				}
				input.Health = &h
			}

			resp, err := linear.ProjectPostCreate(ctx, client, input)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			u := resp.ProjectUpdateCreate.ProjectUpdate
			output.PrintJSON(map[string]any{
				"created":   resp.ProjectUpdateCreate.Success,
				"id":        u.Id,
				"url":       u.Url,
				"health":    string(u.Health),
				"createdAt": u.CreatedAt,
			})
		},
	}

	cmd.Flags().StringVar(&health, "health", "", "Project health: on-track|at-risk|off-track")
	parent.AddCommand(cmd)
}
