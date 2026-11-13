package resolvers

import (
	"context"
	"strings"

	apierrors "github.com/shhac/lin/internal/errors"
)

func ctx() context.Context { return context.Background() }

// requireInput rejects blank resolver input before any lookup runs. Linear
// treats an empty eqIgnoreCase comparator as match-all, so a blank name would
// otherwise resolve to whichever entity the API happens to return first.
func requireInput(noun, input string) error {
	if strings.TrimSpace(input) != "" {
		return nil
	}
	return emptyInputErr(noun)
}

func emptyInputErr(noun string) error {
	return apierrors.Newf(apierrors.FixableByAgent, "%s must not be empty", noun)
}

// matchByNameOrID returns the item whose id equals input, else every item
// whose name matches input case-insensitively: one match resolves, none is
// not-found, and several are ambiguous.
func matchByNameOrID[T any](input string, items []T, id, name func(T) string) []T {
	for _, it := range items {
		if id(it) == input {
			return []T{it}
		}
	}
	lower := strings.ToLower(input)
	var matches []T
	for _, it := range items {
		if strings.ToLower(name(it)) == lower {
			matches = append(matches, it)
		}
	}
	return matches
}

// labelNotFoundErr builds the shared "<noun> not found" error for the label
// resolvers; noun is "label" or "project label" so the message stays exact.
func labelNotFoundErr(noun, input string, names []string) error {
	return apierrors.Newf(apierrors.FixableByAgent, "%s not found: %q, available labels: %s", noun, input, formatChoices(names))
}

// nameMatchLimit caps a name-fallback list query. One result is never enough:
// Linear allows duplicate names, so the resolver must see a second match to
// refuse rather than silently pick whichever the API returned first.
const nameMatchLimit = 10

// ambiguousErr builds the shared "ambiguous <noun>" error, listing each match
// as describe renders it; fix tells the agent how to pick one.
func ambiguousErr[T any](noun, plural, input string, matches []T, describe func(T) string, fix string) error {
	parts := make([]string, len(matches))
	for i, m := range matches {
		parts[i] = describe(m)
	}
	return apierrors.Newf(apierrors.FixableByAgent, "ambiguous %s: %q matches %d %s: %s, %s", noun, input, len(parts), plural, formatChoices(parts), fix)
}

// ambiguousLabelErr is ambiguousErr for the label resolvers; hint is an
// optional trailing suffix (the issue resolver's --team tip).
func ambiguousLabelErr[T any](noun, input string, matches []T, describe func(T) string, hint string) error {
	return ambiguousErr(noun, "labels", input, matches, describe, "use the label ID to disambiguate"+hint)
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
