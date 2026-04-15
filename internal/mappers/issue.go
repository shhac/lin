package mappers

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
