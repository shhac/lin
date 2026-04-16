package user

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerMe(user *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "me",
		Short: "Current authenticated user",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			client := linear.GetClient()

			resp, err := linear.Viewer(context.Background(), client)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			v := resp.Viewer
			output.PrintJSON(map[string]any{
				"id":          v.Id,
				"name":        v.Name,
				"email":       v.Email,
				"displayName": v.DisplayName,
				"organization": map[string]any{
					"id":   v.Organization.Id,
					"name": v.Organization.Name,
				},
			})
		},
	}
	user.AddCommand(cmd)
}
