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
		teamLabels, err := linear.FetchAll(func(first int, after *string) ([]linear.TeamLabelsTeamLabelsIssueLabelConnectionNodesIssueLabel, bool, *string, error) {
			resp, err := linear.TeamLabels(ctx(), client, teamID, first, after)
			if err != nil {
				return nil, false, nil, err
			}
			return resp.Team.Labels.Nodes, resp.Team.Labels.PageInfo.HasNextPage, resp.Team.Labels.PageInfo.EndCursor, nil
		})
		if err != nil {
			return nil, err
		}
		allLabels, err := linear.FetchAll(func(first int, after *string) ([]linear.LabelListIssueLabelsIssueLabelConnectionNodesIssueLabel, bool, *string, error) {
			resp, err := linear.LabelList(ctx(), client, first, after)
			if err != nil {
				return nil, false, nil, err
			}
			return resp.IssueLabels.Nodes, resp.IssueLabels.PageInfo.HasNextPage, resp.IssueLabels.PageInfo.EndCursor, nil
		})
		if err != nil {
			return nil, err
		}
		teamIDs := map[string]bool{}
		var labels []ResolvedLabel
		for _, l := range teamLabels {
			teamIDs[l.Id] = true
			labels = append(labels, ResolvedLabel{ID: l.Id, Name: l.Name})
		}
		for _, l := range allLabels {
			if !teamIDs[l.Id] {
				labels = append(labels, ResolvedLabel{ID: l.Id, Name: l.Name})
			}
		}
		return labels, nil
	}
	allLabels, err := linear.FetchAll(func(first int, after *string) ([]linear.LabelListIssueLabelsIssueLabelConnectionNodesIssueLabel, bool, *string, error) {
		resp, err := linear.LabelList(ctx(), client, first, after)
		if err != nil {
			return nil, false, nil, err
		}
		return resp.IssueLabels.Nodes, resp.IssueLabels.PageInfo.HasNextPage, resp.IssueLabels.PageInfo.EndCursor, nil
	})
	if err != nil {
		return nil, err
	}
	labels := make([]ResolvedLabel, len(allLabels))
	for i, l := range allLabels {
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
		return ResolvedLabel{}, fmt.Errorf("label not found: %q, available labels: %s", input, strings.Join(names, ", "))
	}
	var ambiguous []string
	for _, l := range matches {
		ambiguous = append(ambiguous, fmt.Sprintf("%s (%s)", l.Name, l.ID))
	}
	hint := ""
	if !teamScoped {
		hint = ", tip: use --team to narrow scope"
	}
	return ResolvedLabel{}, fmt.Errorf("ambiguous label: %q matches %d labels: %s, use the label ID to disambiguate%s", input, len(matches), strings.Join(ambiguous, ", "), hint)
}
