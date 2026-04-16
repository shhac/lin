package linear

import (
	"net/http"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/credential"
	"github.com/shhac/lin/internal/output"
)

type authTransport struct {
	apiKey  string
	wrapped http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.apiKey)
	return t.wrapped.RoundTrip(req)
}

func mustResolveAPIKey() string {
	apiKey := credential.Resolve()
	if apiKey == "" {
		output.PrintError("Not authenticated. Run `lin auth login` to connect your Linear workspace.")
	}
	return apiKey
}

// GetClient returns a genqlient GraphQL client authenticated with the
// resolved API key. Exits with a JSON error if no key is available.
func GetClient() graphql.Client {
	return graphql.NewClient(defaultAPIURL, &http.Client{
		Transport: &authTransport{apiKey: mustResolveAPIKey(), wrapped: http.DefaultTransport},
	})
}

// GetRawClient returns our custom Client for raw GraphQL queries.
func GetRawClient() *Client {
	return NewClient(mustResolveAPIKey())
}

// GetAPIKey returns the resolved API key, exiting if unavailable.
func GetAPIKey() string {
	return mustResolveAPIKey()
}
