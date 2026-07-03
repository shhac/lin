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
	apiKey, err := credential.ResolveForClient()
	if err != nil {
		output.WriteError(err)
	}
	if apiKey == "" {
		output.PrintError("Not authenticated. Run `lin auth login` to connect your Linear workspace.")
	}
	return apiKey
}

// ClientWithKey returns a GraphQL client authenticated with an explicit API key
// against the configured endpoint. Used by flows that must validate a key
// before it is stored — MCP credential enrollment — where the key is not yet
// resolvable from the store. Honors the linear.Configure base-URL override so
// tests can point it at a fake endpoint.
func ClientWithKey(apiKey string) graphql.Client {
	return graphql.NewClient(endpointURL(), httpClient(apiKey))
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
