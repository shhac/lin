package errors

import (
	stderrors "errors"
	"regexp"
	"strings"

	"github.com/Khan/genqlient/graphql"
)

// genqlient's gqlerror formats as "input:3: fieldName Actual error message"
// The "input:N:" prefix refers to the line in the GraphQL query — useless to
// CLI users who can't edit the query.
var gqlerrorPrefixRE = regexp.MustCompile(`^input:\d+:\s*\w+\s+`)

// ClassifyGraphQLError extracts a clean, actionable error from a genqlient error.
func ClassifyGraphQLError(err error) *APIError {
	if err == nil {
		return nil
	}

	// Check for HTTP-level errors (auth, rate limit, server errors)
	var httpErr *graphql.HTTPError
	if stderrors.As(err, &httpErr) {
		return classifyHTTPError(httpErr)
	}

	// Clean up GraphQL error messages
	msg := cleanGraphQLError(err.Error())

	// Classify by message content
	lower := strings.ToLower(msg)

	if strings.Contains(lower, "not found") || strings.Contains(lower, "entity not found") {
		return New("not found: "+extractEntity(msg), FixableByAgent).
			WithHint("check the ID or key — use 'list' to see available items")
	}

	if strings.Contains(lower, "authentication") || strings.Contains(lower, "unauthorized") {
		return New("authentication failed: "+msg, FixableByHuman).
			WithHint("check your API key with 'lin auth status'")
	}

	if strings.Contains(lower, "forbidden") || strings.Contains(lower, "permission") {
		return New("permission denied: "+msg, FixableByHuman).
			WithHint("your API key may not have sufficient permissions")
	}

	return Wrap(stderrors.New(msg), FixableByAgent)
}

func classifyHTTPError(httpErr *graphql.HTTPError) *APIError {
	msg := extractHTTPErrorMessage(httpErr)

	switch {
	case httpErr.StatusCode == 401:
		return New("authentication failed: "+msg, FixableByHuman).
			WithHint("check your API key with 'lin auth status'")
	case httpErr.StatusCode == 403:
		return New("permission denied: "+msg, FixableByHuman).
			WithHint("your API key may not have sufficient permissions")
	case httpErr.StatusCode == 429:
		return New("rate limited", FixableByRetry).
			WithHint("Linear rate limit hit — wait and retry")
	case httpErr.StatusCode >= 500:
		return New("Linear API error: "+msg, FixableByRetry).
			WithHint("server error — retry in a few seconds")
	default:
		return New(msg, FixableByAgent)
	}
}

func extractHTTPErrorMessage(httpErr *graphql.HTTPError) string {
	if len(httpErr.Response.Errors) > 0 {
		return cleanGraphQLError(httpErr.Response.Errors[0].Message)
	}
	return strings.TrimSpace(httpErr.Error())
}

// cleanGraphQLError strips genqlient's "input:N: fieldName " prefix and
// trailing newlines from error messages.
func cleanGraphQLError(msg string) string {
	msg = strings.TrimSpace(msg)
	// Strip "input:N: fieldName " prefix
	msg = gqlerrorPrefixRE.ReplaceAllString(msg, "")
	// Strip "returned error N:" prefix from HTTP errors
	if strings.HasPrefix(msg, "returned error") {
		if idx := strings.Index(msg, ":"); idx >= 0 {
			msg = strings.TrimSpace(msg[idx+1:])
		}
	}
	return msg
}

// extractEntity tries to pull out the entity type from a "not found" error.
// e.g., "Entity not found: Issue" → "Issue"
// e.g., "Could not find referenced Issue." → "Issue"
func extractEntity(msg string) string {
	if idx := strings.LastIndex(msg, ": "); idx >= 0 {
		entity := strings.TrimRight(msg[idx+2:], ".")
		if entity != "" {
			return entity
		}
	}
	if strings.Contains(msg, "Issue") {
		return "Issue"
	}
	if strings.Contains(msg, "Project") {
		return "Project"
	}
	if strings.Contains(msg, "Document") {
		return "Document"
	}
	return "entity"
}
