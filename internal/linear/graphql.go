package linear

import (
	"fmt"
	"net/http"
	"time"

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

type debugTransport struct {
	wrapped http.RoundTripper
}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	_, _ = fmt.Fprintf(output.Stderr(), "debug: %s %s\n", req.Method, req.URL.Redacted())
	return t.wrapped.RoundTrip(req)
}

type Options struct {
	BaseURL   string
	TimeoutMS int
	Debug     bool
}

var clientOptions Options

func Configure(opts Options) {
	clientOptions = opts
}

func endpointURL() string {
	if clientOptions.BaseURL != "" {
		return clientOptions.BaseURL
	}
	return defaultAPIURL
}

func httpClient(apiKey string) *http.Client {
	transport := http.RoundTripper(&authTransport{apiKey: apiKey, wrapped: http.DefaultTransport})
	if clientOptions.Debug {
		transport = &debugTransport{wrapped: transport}
	}
	client := &http.Client{Transport: transport}
	if clientOptions.TimeoutMS > 0 {
		client.Timeout = time.Duration(clientOptions.TimeoutMS) * time.Millisecond
	}
	return client
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
	apiKey := mustResolveAPIKey()
	return graphql.NewClient(endpointURL(), httpClient(apiKey))
}

// GetRawClient returns our custom Client for raw GraphQL queries.
func GetRawClient() *Client {
	apiKey := mustResolveAPIKey()
	client := NewClientWithHTTP(apiKey, httpClient(apiKey))
	client.baseURL = endpointURL()
	return client
}

// GetAPIKey returns the resolved API key, exiting if unavailable.
func GetAPIKey() string {
	return mustResolveAPIKey()
}
