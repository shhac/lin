package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"

	"github.com/Khan/genqlient/graphql"
)

var opNameRE = regexp.MustCompile(`(?:query|mutation)\s+(\w+)`)

// MockLinear creates a test HTTP server that responds to GraphQL operations.
// Operations are matched by extracting the operation name from the query.
type MockLinear struct {
	Server   *httptest.Server
	mu       sync.RWMutex
	handlers map[string]any
}

type graphqlRequest struct {
	Query         string         `json:"query"`
	Variables     map[string]any `json:"variables,omitempty"`
	OperationName string         `json:"operationName,omitempty"`
}

func NewMockLinear() *MockLinear {
	m := &MockLinear{
		handlers: make(map[string]any),
	}
	m.Server = httptest.NewServer(http.HandlerFunc(m.handle))
	return m
}

// Handle registers a response for a given GraphQL operation name.
func (m *MockLinear) Handle(operationName string, response any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[operationName] = response
}

// Client returns a genqlient graphql.Client pointing at the mock server.
func (m *MockLinear) Client() graphql.Client {
	return graphql.NewClient(m.Server.URL, m.Server.Client())
}

// Close shuts down the test server.
func (m *MockLinear) Close() {
	m.Server.Close()
}

func (m *MockLinear) handle(w http.ResponseWriter, r *http.Request) {
	var req graphqlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeGraphQLError(w, fmt.Sprintf("failed to decode request: %v", err))
		return
	}

	opName := req.OperationName
	if opName == "" {
		if matches := opNameRE.FindStringSubmatch(req.Query); len(matches) > 1 {
			opName = matches[1]
		}
	}

	if opName == "" {
		writeGraphQLError(w, "could not determine operation name from query")
		return
	}

	m.mu.RLock()
	handler, ok := m.handlers[opName]
	m.mu.RUnlock()

	if !ok {
		writeGraphQLError(w, fmt.Sprintf("no handler for operation: %s", opName))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"data": handler})
}

func writeGraphQLError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"errors": []map[string]string{{"message": msg}},
	})
}
