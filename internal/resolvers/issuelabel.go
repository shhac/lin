package resolvers

import (
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
)

type ResolvedIssueLabel struct {
	ID      string
	Name    string
	TeamKey string // empty for workspace-wide labels
}

func ResolveIssueLabels(client graphql.Client, input, teamID string) ([]string, error) {
	inputs := splitAndTrim(input)
	labels, err := fetchIssueLabels(client, teamID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(inputs))
	for _, raw := range inputs {
		l, err := resolveOneIssueLabel(raw, labels, teamID != "")
		if err != nil {
			return nil, err
		}
		ids = append(ids, l.ID)
	}
	return ids, nil
}

func fetchOrgIssueLabels(client graphql.Client) ([]ResolvedIssueLabel, error) {
	nodes, err := linear.FetchAll(func(first int, after *string) ([]linear.IssueLabelFields, bool, *string, error) {
		resp, err := linear.IssueLabelList(ctx(), client, first, after, nil)
		if err != nil {
			return nil, false, nil, err
		}
		out := make([]linear.IssueLabelFields, len(resp.IssueLabels.Nodes))
		for i, n := range resp.IssueLabels.Nodes {
			out[i] = n.IssueLabelFields
		}
		return out, resp.IssueLabels.PageInfo.HasNextPage, resp.IssueLabels.PageInfo.EndCursor, nil
	})
	if err != nil {
		return nil, err
	}
	labels := make([]ResolvedIssueLabel, len(nodes))
	for i, l := range nodes {
		labels[i] = ResolvedIssueLabel{ID: l.Id, Name: l.Name, TeamKey: teamKeyOfIssueLabel(l)}
	}
	return labels, nil
}

func fetchIssueLabels(client graphql.Client, teamID string) ([]ResolvedIssueLabel, error) {
	if teamID == "" {
		return fetchOrgIssueLabels(client)
	}

	teamNodes, err := linear.FetchAll(func(first int, after *string) ([]linear.TeamLabelsTeamLabelsIssueLabelConnectionNodesIssueLabel, bool, *string, error) {
		resp, err := linear.TeamLabels(ctx(), client, teamID, first, after)
		if err != nil {
			return nil, false, nil, err
		}
		return resp.Team.Labels.Nodes, resp.Team.Labels.PageInfo.HasNextPage, resp.Team.Labels.PageInfo.EndCursor, nil
	})
	if err != nil {
		return nil, err
	}
	orgLabels, err := fetchOrgIssueLabels(client)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	labels := make([]ResolvedIssueLabel, 0, len(teamNodes)+len(orgLabels))
	for _, n := range teamNodes {
		seen[n.Id] = true
		labels = append(labels, ResolvedIssueLabel{ID: n.Id, Name: n.Name, TeamKey: teamKeyOfIssueLabel(n.IssueLabelFields)})
	}
	for _, l := range orgLabels {
		if !seen[l.ID] {
			labels = append(labels, l)
		}
	}
	return labels, nil
}

func teamKeyOfIssueLabel(l linear.IssueLabelFields) string {
	if l.Team == nil {
		return ""
	}
	return l.Team.Key
}

func resolveOneIssueLabel(input string, labels []ResolvedIssueLabel, teamScoped bool) (ResolvedIssueLabel, error) {
	for _, l := range labels {
		if l.ID == input {
			return l, nil
		}
	}
	lower := strings.ToLower(input)
	var matches []ResolvedIssueLabel
	for _, l := range labels {
		if strings.ToLower(l.Name) == lower {
			matches = append(matches, l)
		}
	}

	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) == 0 {
		return ResolvedIssueLabel{}, issueLabelNotFoundErr(input, labels)
	}
	return ResolvedIssueLabel{}, ambiguousIssueLabelErr(input, matches, teamScoped)
}

func issueLabelNotFoundErr(input string, labels []ResolvedIssueLabel) error {
	names := make([]string, len(labels))
	for i, l := range labels {
		names[i] = l.Name
	}
	return fmt.Errorf("label not found: %q, available labels: %s", input, strings.Join(names, ", "))
}

func ambiguousIssueLabelErr(input string, matches []ResolvedIssueLabel, teamScoped bool) error {
	parts := make([]string, len(matches))
	for i, l := range matches {
		if l.TeamKey != "" {
			parts[i] = fmt.Sprintf("%s (id: %s, team: %s)", l.Name, l.ID, l.TeamKey)
		} else {
			parts[i] = fmt.Sprintf("%s (id: %s, workspace)", l.Name, l.ID)
		}
	}
	hint := ""
	if !teamScoped {
		hint = ", tip: use --team to narrow scope"
	}
	return fmt.Errorf("ambiguous label: %q matches %d labels: %s, use the label ID to disambiguate%s", input, len(matches), strings.Join(parts, ", "), hint)
}
