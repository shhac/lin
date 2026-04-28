package resolvers

import (
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
)

type ResolvedProjectLabel struct {
	ID   string
	Name string
}

// ResolveProjectLabels translates a comma-separated list of project label
// names or UUIDs into label IDs. Project labels are workspace-scoped; there
// is no team scope to consider.
func ResolveProjectLabels(client graphql.Client, input string) ([]string, error) {
	inputs := splitAndTrim(input)
	labels, err := fetchProjectLabels(client)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(inputs))
	for _, raw := range inputs {
		l, err := resolveOneProjectLabel(raw, labels)
		if err != nil {
			return nil, err
		}
		ids = append(ids, l.ID)
	}
	return ids, nil
}

func fetchProjectLabels(client graphql.Client) ([]ResolvedProjectLabel, error) {
	nodes, err := linear.FetchAll(func(first int, after *string) ([]linear.ProjectLabelFields, bool, *string, error) {
		resp, err := linear.ProjectLabelList(ctx(), client, first, after, nil)
		if err != nil {
			return nil, false, nil, err
		}
		out := make([]linear.ProjectLabelFields, len(resp.ProjectLabels.Nodes))
		for i, n := range resp.ProjectLabels.Nodes {
			out[i] = n.ProjectLabelFields
		}
		return out, resp.ProjectLabels.PageInfo.HasNextPage, resp.ProjectLabels.PageInfo.EndCursor, nil
	})
	if err != nil {
		return nil, err
	}
	labels := make([]ResolvedProjectLabel, len(nodes))
	for i, l := range nodes {
		labels[i] = ResolvedProjectLabel{ID: l.Id, Name: l.Name}
	}
	return labels, nil
}

func resolveOneProjectLabel(input string, labels []ResolvedProjectLabel) (ResolvedProjectLabel, error) {
	for _, l := range labels {
		if l.ID == input {
			return l, nil
		}
	}
	lower := strings.ToLower(input)
	var matches []ResolvedProjectLabel
	for _, l := range labels {
		if strings.ToLower(l.Name) == lower {
			matches = append(matches, l)
		}
	}

	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) == 0 {
		return ResolvedProjectLabel{}, projectLabelNotFoundErr(input, labels)
	}
	return ResolvedProjectLabel{}, ambiguousProjectLabelErr(input, matches)
}

func projectLabelNotFoundErr(input string, labels []ResolvedProjectLabel) error {
	names := make([]string, len(labels))
	for i, l := range labels {
		names[i] = l.Name
	}
	return fmt.Errorf("project label not found: %q, available labels: %s", input, strings.Join(names, ", "))
}

func ambiguousProjectLabelErr(input string, matches []ResolvedProjectLabel) error {
	parts := make([]string, len(matches))
	for i, l := range matches {
		parts[i] = fmt.Sprintf("%s (id: %s)", l.Name, l.ID)
	}
	return fmt.Errorf("ambiguous project label: %q matches %d labels: %s, use the label ID to disambiguate", input, len(matches), strings.Join(parts, ", "))
}
