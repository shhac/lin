package resolvers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
)

func ctx() context.Context { return context.Background() }

type ResolvedUser struct {
	ID          string
	Name        string
	Email       string
	DisplayName string
}

func ResolveUser(client graphql.Client, input string) (ResolvedUser, error) {
	resp, err := linear.UserList(ctx(), client, nil, 250, nil)
	if err != nil {
		return ResolvedUser{}, err
	}
	lower := strings.ToLower(input)
	var matches []ResolvedUser
	for _, u := range resp.Users.Nodes {
		if u.Id == input ||
			strings.ToLower(u.Name) == lower ||
			strings.ToLower(u.Email) == lower ||
			strings.ToLower(u.DisplayName) == lower {
			matches = append(matches, ResolvedUser{
				ID: u.Id, Name: u.Name, Email: u.Email, DisplayName: u.DisplayName,
			})
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) == 0 {
		var names []string
		for _, u := range resp.Users.Nodes {
			names = append(names, fmt.Sprintf("%s <%s>", u.Name, u.Email))
		}
		return ResolvedUser{}, fmt.Errorf("User not found: %q. Known users: %s", input, strings.Join(names, ", "))
	}
	var ambiguous []string
	for _, u := range matches {
		ambiguous = append(ambiguous, fmt.Sprintf("%s <%s> (%s)", u.Name, u.Email, u.ID))
	}
	return ResolvedUser{}, fmt.Errorf("Ambiguous user: %q matches %d users: %s. Use a unique name, email, or ID.", input, len(matches), strings.Join(ambiguous, ", "))
}

type ResolvedProject struct {
	ID     string
	Name   string
	SlugId string
}

func ResolveProject(client graphql.Client, input string) (ResolvedProject, error) {
	resp, err := linear.ProjectGet(ctx(), client, input)
	if err == nil {
		return ResolvedProject{
			ID: resp.Project.Id, Name: resp.Project.Name, SlugId: resp.Project.SlugId,
		}, nil
	}
	filter := filters.BuildProjectFilter(input)
	listResp, err := linear.ProjectList(ctx(), client, filter, 1, nil)
	if err != nil {
		return ResolvedProject{}, err
	}
	if len(listResp.Projects.Nodes) == 0 {
		return ResolvedProject{}, fmt.Errorf("Project not found: %q. Provide a UUID, slug ID, or exact name.", input)
	}
	p := listResp.Projects.Nodes[0]
	return ResolvedProject{ID: p.Id, Name: p.Name, SlugId: p.SlugId}, nil
}

type ResolvedTeam struct {
	ID   string
	Name string
	Key  string
}

func ResolveTeam(client graphql.Client, input string) (ResolvedTeam, error) {
	resp, err := linear.TeamGet(ctx(), client, input)
	if err == nil {
		return ResolvedTeam{ID: resp.Team.Id, Name: resp.Team.Name, Key: resp.Team.Key}, nil
	}
	filter := filters.BuildTeamFilter(input)
	listResp, err := linear.TeamList(ctx(), client, filter, 50, nil)
	if err != nil {
		return ResolvedTeam{}, err
	}
	if len(listResp.Teams.Nodes) > 0 {
		t := listResp.Teams.Nodes[0]
		return ResolvedTeam{ID: t.Id, Name: t.Name, Key: t.Key}, nil
	}
	allResp, err := linear.TeamList(ctx(), client, nil, 250, nil)
	if err != nil {
		return ResolvedTeam{}, err
	}
	var keys []string
	for _, t := range allResp.Teams.Nodes {
		keys = append(keys, fmt.Sprintf("%s (%s)", t.Key, t.Name))
	}
	hint := "none"
	if len(keys) > 0 {
		hint = strings.Join(keys, ", ")
	}
	return ResolvedTeam{}, fmt.Errorf("Team not found: %q. Known teams: %s. Provide a UUID, key, or exact name.", input, hint)
}

type ResolvedDocument struct {
	ID     string
	SlugId string
	Title  string
}

func ResolveDocument(client graphql.Client, input string) (ResolvedDocument, error) {
	resp, err := linear.DocumentGet(ctx(), client, input)
	if err == nil {
		return ResolvedDocument{
			ID: resp.Document.Id, SlugId: resp.Document.SlugId, Title: resp.Document.Title,
		}, nil
	}
	filter := &linear.DocumentFilter{
		SlugId: &linear.StringComparator{Eq: strPtr(input)},
	}
	listResp, err := linear.DocumentList(ctx(), client, filter, 1, nil, nil)
	if err != nil {
		return ResolvedDocument{}, err
	}
	if len(listResp.Documents.Nodes) == 0 {
		return ResolvedDocument{}, fmt.Errorf("Document not found: %q. Provide a UUID or slug ID.", input)
	}
	d := listResp.Documents.Nodes[0]
	return ResolvedDocument{ID: d.Id, SlugId: d.SlugId, Title: d.Title}, nil
}

type ResolvedWorkflowState struct {
	ID   string
	Name string
}

func ResolveWorkflowState(client graphql.Client, name, teamID string) (ResolvedWorkflowState, error) {
	filter := &linear.WorkflowStateFilter{
		Team: &linear.TeamFilter{Id: &linear.IDComparator{Eq: strPtr(teamID)}},
	}
	resp, err := linear.WorkflowStates(ctx(), client, filter)
	if err != nil {
		return ResolvedWorkflowState{}, err
	}
	lower := strings.ToLower(name)
	for _, s := range resp.WorkflowStates.Nodes {
		if strings.ToLower(s.Name) == lower {
			return ResolvedWorkflowState{ID: s.Id, Name: s.Name}, nil
		}
	}
	seen := map[string]bool{}
	var validNames []string
	for _, s := range resp.WorkflowStates.Nodes {
		if !seen[s.Name] {
			seen[s.Name] = true
			validNames = append(validNames, s.Name)
		}
	}
	return ResolvedWorkflowState{}, fmt.Errorf("Unknown status: %q. Valid values: %s", name, strings.Join(validNames, " | "))
}

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

type ResolvedRoadmap struct {
	ID     string
	Name   string
	SlugId string
}

func ResolveRoadmap(client graphql.Client, input string) (ResolvedRoadmap, error) {
	resp, err := linear.RoadmapGet(ctx(), client, input)
	if err == nil {
		return ResolvedRoadmap{
			ID: resp.Roadmap.Id, Name: resp.Roadmap.Name, SlugId: resp.Roadmap.SlugId,
		}, nil
	}
	listResp, err := linear.RoadmapList(ctx(), client, 250, nil)
	if err != nil {
		return ResolvedRoadmap{}, err
	}
	lower := strings.ToLower(input)
	for _, r := range listResp.Roadmaps.Nodes {
		if r.SlugId == input || strings.ToLower(r.Name) == lower {
			return ResolvedRoadmap{ID: r.Id, Name: r.Name, SlugId: r.SlugId}, nil
		}
	}
	var names []string
	for _, r := range listResp.Roadmaps.Nodes {
		names = append(names, fmt.Sprintf("%s (%s)", r.Name, r.SlugId))
	}
	hint := "none"
	if len(names) > 0 {
		hint = strings.Join(names, ", ")
	}
	return ResolvedRoadmap{}, fmt.Errorf("Roadmap not found: %q. Known roadmaps: %s. Provide a UUID, slug ID, or exact name.", input, hint)
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

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func strPtr(s string) *string { return &s }
