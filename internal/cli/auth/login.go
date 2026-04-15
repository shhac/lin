package auth

import (
	"context"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

type loginTransport struct {
	apiKey string
}

func (t *loginTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.apiKey)
	return http.DefaultTransport.RoundTrip(req)
}

func registerLogin(auth *cobra.Command) {
	var alias string

	cmd := &cobra.Command{
		Use:   "login <api-key>",
		Short: "Store Linear API key (auto-detects workspace)",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			apiKey := args[0]

			httpClient := &http.Client{Transport: &loginTransport{apiKey: apiKey}}
			client := graphql.NewClient("https://api.linear.app/graphql", httpClient)

			ctx := context.Background()
			resp, err := linear.Viewer(ctx, client)
			if err != nil {
				output.PrintError(err.Error())
				return
			}

			viewer := resp.Viewer
			org := viewer.Organization
			wsAlias := alias
			if wsAlias == "" {
				wsAlias = org.UrlKey
			}

			if err := config.StoreLogin(wsAlias, config.Workspace{
				APIKey: apiKey,
				Name:   org.Name,
				URLKey: org.UrlKey,
			}); err != nil {
				output.PrintError(err.Error())
				return
			}

			isDefault := config.GetDefaultWorkspace() == wsAlias
			output.PrintJSON(map[string]any{
				"ok": true,
				"user": map[string]any{
					"id":    viewer.Id,
					"name":  viewer.Name,
					"email": viewer.Email,
				},
				"workspace": map[string]any{
					"alias":   wsAlias,
					"name":    org.Name,
					"urlKey":  org.UrlKey,
					"default": isDefault,
				},
				"hint": "To add another workspace, run: lin auth login <other-api-key>",
			})
		},
	}

	cmd.Flags().StringVar(&alias, "alias", "", "Custom workspace alias (default: org urlKey)")
	auth.AddCommand(cmd)
}
