package resolvers

import (
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
)

type ResolvedLabel struct {
	ID   string
	Name string
}

func ResolveLabels(client graphql.Client, input, teamID string) ([]string, error) {
	inputs := splitAndTrim(input)
	labels, err := fetchLabels(client, teamID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(inputs))
	for _, raw := range inputs {
		l, err := resolveOneLabel(raw, labels, teamID != "")
		if err != nil {
			return nil, err
		}
		ids = append(ids, l.ID)
	}
	return ids, nil
}

func fetchLabels(client graphql.Client, teamID string) ([]ResolvedLabel, error) {
	if teamID != "" {
		teamResp, err := linear.TeamLabels(ctx(), client, teamID, 250, nil)
		if err != nil {
			return nil, err
		}
		allResp, err := linear.LabelList(ctx(), client, 250, nil)
		if err != nil {
			return nil, err
		}
		teamIDs := map[string]bool{}
		var labels []ResolvedLabel
		for _, l := range teamResp.Team.Labels.Nodes {
			teamIDs[l.Id] = true
			labels = append(labels, ResolvedLabel{ID: l.Id, Name: l.Name})
		}
		for _, l := range allResp.IssueLabels.Nodes {
			if !teamIDs[l.Id] {
				labels = append(labels, ResolvedLabel{ID: l.Id, Name: l.Name})
			}
		}
		return labels, nil
	}
	resp, err := linear.LabelList(ctx(), client, 250, nil)
	if err != nil {
		return nil, err
	}
	labels := make([]ResolvedLabel, len(resp.IssueLabels.Nodes))
	for i, l := range resp.IssueLabels.Nodes {
		labels[i] = ResolvedLabel{ID: l.Id, Name: l.Name}
	}
	return labels, nil
}

func resolveOneLabel(input string, labels []ResolvedLabel, teamScoped bool) (ResolvedLabel, error) {
	for _, l := range labels {
		if l.ID == input {
			return l, nil
		}
	}
	lower := strings.ToLower(input)
	var matches []ResolvedLabel
	for _, l := range labels {
		if strings.ToLower(l.Name) == lower {
			matches = append(matches, l)
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) == 0 {
		var names []string
		for _, l := range labels {
			names = append(names, l.Name)
		}
		return ResolvedLabel{}, fmt.Errorf("Label not found: %q. Available labels: %s", input, strings.Join(names, ", "))
	}
	var ambiguous []string
	for _, l := range matches {
		ambiguous = append(ambiguous, fmt.Sprintf("%s (%s)", l.Name, l.ID))
	}
	hint := ""
	if !teamScoped {
		hint = " Tip: use --team to narrow scope."
	}
	return ResolvedLabel{}, fmt.Errorf("Ambiguous label: %q matches %d labels: %s. Use the label ID to disambiguate.%s", input, len(matches), strings.Join(ambiguous, ", "), hint)
}
