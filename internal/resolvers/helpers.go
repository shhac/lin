package resolvers

import (
	"context"
	"strings"
)

func ctx() context.Context { return context.Background() }

// formatChoices renders a list of valid options for an error message as a
// comma-separated string, the shape resolver "not found / ambiguous" errors use
// to tell the agent what it could have passed instead.
func formatChoices(choices []string) string {
	return strings.Join(choices, ", ")
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
