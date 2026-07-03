package resolvers

import (
	"context"
	"strings"

	apierrors "github.com/shhac/lin/internal/errors"
)

func ctx() context.Context { return context.Background() }

// matchKind is the outcome of matchByNameOrID.
type matchKind int

const (
	matchFound matchKind = iota
	matchNone
	matchAmbiguous
)

// matchByNameOrID finds the item whose id equals input, else the item(s) whose
// name matches input case-insensitively. It returns the single match (kind
// matchFound), the ambiguous set (matchAmbiguous), or nothing (matchNone) — the
// shared resolution shape for name-or-UUID lookups over an already-fetched set.
func matchByNameOrID[T any](input string, items []T, id, name func(T) string) (match T, ambiguous []T, kind matchKind) {
	for _, it := range items {
		if id(it) == input {
			return it, nil, matchFound
		}
	}
	lower := strings.ToLower(input)
	var matches []T
	for _, it := range items {
		if strings.ToLower(name(it)) == lower {
			matches = append(matches, it)
		}
	}
	switch len(matches) {
	case 1:
		return matches[0], nil, matchFound
	case 0:
		var zero T
		return zero, nil, matchNone
	default:
		var zero T
		return zero, matches, matchAmbiguous
	}
}

// labelNotFoundErr builds the shared "<noun> not found" error for the label
// resolvers; noun is "label" or "project label" so the message stays exact.
func labelNotFoundErr(noun, input string, names []string) error {
	return apierrors.Newf(apierrors.FixableByAgent, "%s not found: %q, available labels: %s", noun, input, formatChoices(names))
}

// ambiguousLabelErr builds the shared "ambiguous <noun>" error. parts are the
// pre-formatted match descriptors (the label resolvers differ in how they
// render a match, so the caller formats them); hint is an optional trailing
// suffix (the issue resolver's --team tip).
func ambiguousLabelErr(noun, input string, parts []string, hint string) error {
	return apierrors.Newf(apierrors.FixableByAgent, "ambiguous %s: %q matches %d labels: %s, use the label ID to disambiguate%s", noun, input, len(parts), formatChoices(parts), hint)
}

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
