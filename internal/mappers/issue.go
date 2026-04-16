package mappers

import "github.com/shhac/lin/internal/linear"

// IssueSummaryInput holds the fields needed to produce an issue summary.
// Callers populate this from their specific genqlient response types.
type IssueSummaryInput struct {
	ID            string
	Identifier    string
	Title         string
	BranchName    string
	Priority      float64
	PriorityLabel string
	StateName     string
	StateType     string
	AssigneeID    string // empty if unassigned
	AssigneeName  string // empty if unassigned
	TeamKey       string
}

// FromIssueSummaryFields converts the genqlient IssueSummaryFields fragment
// (used by IssueList, ProjectIssues, CycleGet) into an IssueSummaryInput.
func FromIssueSummaryFields(f linear.IssueSummaryFields) IssueSummaryInput {
	input := IssueSummaryInput{
		ID:            f.Id,
		Identifier:    f.Identifier,
		Title:         f.Title,
		BranchName:    f.BranchName,
		Priority:      f.Priority,
		PriorityLabel: f.PriorityLabel,
		StateName:     f.State.Name,
		StateType:     f.State.Type,
		TeamKey:       f.Team.Key,
	}
	if f.Assignee != nil {
		input.AssigneeID = f.Assignee.Id
		input.AssigneeName = f.Assignee.Name
	}
	return input
}

// FromIssueSearchSummaryFields converts the genqlient IssueSearchSummaryFields
// fragment (used by IssueSearch) into an IssueSummaryInput.
func FromIssueSearchSummaryFields(f linear.IssueSearchSummaryFields) IssueSummaryInput {
	input := IssueSummaryInput{
		ID:            f.Id,
		Identifier:    f.Identifier,
		Title:         f.Title,
		BranchName:    f.BranchName,
		Priority:      f.Priority,
		PriorityLabel: f.PriorityLabel,
		StateName:     f.State.Name,
		StateType:     f.State.Type,
		TeamKey:       f.Team.Key,
	}
	if f.Assignee != nil {
		input.AssigneeID = f.Assignee.Id
		input.AssigneeName = f.Assignee.Name
	}
	return input
}

// MapIssueSummary converts an IssueSummaryInput to the standard output map.
func MapIssueSummary(s IssueSummaryInput) map[string]any {
	m := map[string]any{
		"id":            s.ID,
		"identifier":    s.Identifier,
		"title":         s.Title,
		"branchName":    s.BranchName,
		"priority":      int(s.Priority),
		"priorityLabel": s.PriorityLabel,
		"status":        s.StateName,
		"statusType":    s.StateType,
		"team":          s.TeamKey,
	}
	if s.AssigneeID != "" {
		m["assignee"] = s.AssigneeName
		m["assigneeId"] = s.AssigneeID
	}
	return m
}
