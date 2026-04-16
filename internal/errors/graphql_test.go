package errors

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestCleanGraphQLError(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"input:3: issue Entity not found: Issue\n", "Entity not found: Issue"},
		{"input:5: workflowStates Cannot read properties of null\n", "Cannot read properties of null"},
		{"input:2: viewer Authentication required", "Authentication required"},
		{"plain error message", "plain error message"},
		{"  whitespace around  \n", "whitespace around"},
		{"returned error 401: {\"errors\":[...]}", `{"errors":[...]}`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := cleanGraphQLError(tt.input)
			if got != tt.want {
				t.Errorf("cleanGraphQLError(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractEntity(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Entity not found: Issue", "Issue"},
		{"Entity not found: Project", "Project"},
		{"Could not find referenced Issue.", "Issue"},
		{"not found: Document", "Document"},
		{"something about Issue somewhere", "Issue"},
		{"generic not found", "entity"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractEntity(tt.input)
			if got != tt.want {
				t.Errorf("extractEntity(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestClassifyGraphQLError_NotFound(t *testing.T) {
	err := fmt.Errorf("input:3: issue Entity not found: Issue\n")
	apiErr := ClassifyGraphQLError(err)

	if apiErr.FixableBy != FixableByAgent {
		t.Errorf("FixableBy = %q, want %q", apiErr.FixableBy, FixableByAgent)
	}
	if !strings.Contains(apiErr.Message, "not found") {
		t.Errorf("Message should contain 'not found', got %q", apiErr.Message)
	}
	if !strings.Contains(apiErr.Message, "Issue") {
		t.Errorf("Message should contain entity type 'Issue', got %q", apiErr.Message)
	}
	if apiErr.Hint == "" {
		t.Error("Hint should be set for not-found errors")
	}
}

func TestClassifyGraphQLError_Auth(t *testing.T) {
	err := fmt.Errorf("input:2: viewer Authentication required")
	apiErr := ClassifyGraphQLError(err)

	if apiErr.FixableBy != FixableByHuman {
		t.Errorf("FixableBy = %q, want %q", apiErr.FixableBy, FixableByHuman)
	}
	if !strings.Contains(apiErr.Message, "authentication") {
		t.Errorf("Message should contain 'authentication', got %q", apiErr.Message)
	}
}

func TestClassifyGraphQLError_Generic(t *testing.T) {
	err := fmt.Errorf("something unexpected happened")
	apiErr := ClassifyGraphQLError(err)

	if apiErr.FixableBy != FixableByAgent {
		t.Errorf("FixableBy = %q, want %q", apiErr.FixableBy, FixableByAgent)
	}
	if apiErr.Message != "something unexpected happened" {
		t.Errorf("Message = %q, want original message", apiErr.Message)
	}
}

func TestClassifyGraphQLError_Nil(t *testing.T) {
	if ClassifyGraphQLError(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func makeHTTPError(statusCode int, msg string) *graphql.HTTPError {
	resp := graphql.Response{}
	if msg != "" {
		resp.Errors = gqlerror.List{{Message: msg}}
	}
	return &graphql.HTTPError{StatusCode: statusCode, Response: resp}
}

func TestClassifyHTTPError_401(t *testing.T) {
	apiErr := ClassifyGraphQLError(makeHTTPError(401, "invalid token"))
	if apiErr.FixableBy != FixableByHuman {
		t.Errorf("FixableBy = %q, want %q", apiErr.FixableBy, FixableByHuman)
	}
	if !strings.Contains(apiErr.Message, "authentication") {
		t.Errorf("Message should contain 'authentication', got %q", apiErr.Message)
	}
	if !strings.Contains(apiErr.Hint, "auth") {
		t.Errorf("Hint should mention auth, got %q", apiErr.Hint)
	}
}

func TestClassifyHTTPError_403(t *testing.T) {
	apiErr := ClassifyGraphQLError(makeHTTPError(403, "forbidden"))
	if apiErr.FixableBy != FixableByHuman {
		t.Errorf("FixableBy = %q, want %q", apiErr.FixableBy, FixableByHuman)
	}
	if !strings.Contains(apiErr.Message, "permission") {
		t.Errorf("Message should contain 'permission', got %q", apiErr.Message)
	}
	if !strings.Contains(apiErr.Hint, "permission") {
		t.Errorf("Hint should mention permission, got %q", apiErr.Hint)
	}
}

func TestClassifyHTTPError_429(t *testing.T) {
	apiErr := ClassifyGraphQLError(makeHTTPError(429, ""))
	if apiErr.FixableBy != FixableByRetry {
		t.Errorf("FixableBy = %q, want %q", apiErr.FixableBy, FixableByRetry)
	}
	if !strings.Contains(apiErr.Message, "rate limit") {
		t.Errorf("Message should contain 'rate limit', got %q", apiErr.Message)
	}
	if !strings.Contains(apiErr.Hint, "rate limit") {
		t.Errorf("Hint should mention rate limit, got %q", apiErr.Hint)
	}
}

func TestClassifyHTTPError_500(t *testing.T) {
	apiErr := ClassifyGraphQLError(makeHTTPError(500, "internal server error"))
	if apiErr.FixableBy != FixableByRetry {
		t.Errorf("FixableBy = %q, want %q", apiErr.FixableBy, FixableByRetry)
	}
	if !strings.Contains(apiErr.Hint, "server error") {
		t.Errorf("Hint should mention server error, got %q", apiErr.Hint)
	}
}

func TestClassifyHTTPError_418(t *testing.T) {
	apiErr := ClassifyGraphQLError(makeHTTPError(418, "i'm a teapot"))
	if apiErr.FixableBy != FixableByAgent {
		t.Errorf("FixableBy = %q, want %q", apiErr.FixableBy, FixableByAgent)
	}
}
