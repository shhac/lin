package auth

import (
	"context"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/credential"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerStatus(auth *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current auth state and workspace info",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			apiKey := credential.Resolve()
			if apiKey == "" {
				output.PrintJSON(map[string]any{"authenticated": false})
				return
			}

			httpClient := &http.Client{Transport: &loginTransport{apiKey: apiKey}}
			client := graphql.NewClient("https://api.linear.app/graphql", httpClient)

			ctx := context.Background()
			resp, err := linear.Viewer(ctx, client)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			viewer := resp.Viewer
			org := viewer.Organization
			workspaces := config.GetWorkspaces()
			defaultWs := config.GetDefaultWorkspace()

			var otherWorkspaces []map[string]any
			for alias, ws := range workspaces {
				if alias == defaultWs {
					continue
				}
				otherWorkspaces = append(otherWorkspaces, map[string]any{
					"alias":  alias,
					"name":   ws.Name,
					"urlKey": ws.URLKey,
				})
			}

			result := map[string]any{
				"authenticated": true,
				"source":        credential.Source(),
				"user": map[string]any{
					"id":    viewer.Id,
					"name":  viewer.Name,
					"email": viewer.Email,
				},
				"organization": map[string]any{
					"id":     org.Id,
					"name":   org.Name,
					"urlKey": org.UrlKey,
				},
				"activeWorkspace": defaultWs,
			}
			if len(otherWorkspaces) > 0 {
				result["otherWorkspaces"] = otherWorkspaces
			}

			output.PrintJSON(result)
		},
	}

	auth.AddCommand(cmd)
}
