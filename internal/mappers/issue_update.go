package mappers

import "github.com/shhac/lin/internal/linear"

// MapIssueUpdateResult renders the issue as Linear reports it after an
// update, so a command can echo the value it actually set. Keys and shapes
// follow MapIssueDetail; a key is present only when that field is set.
func MapIssueUpdateResult(issue linear.IssueUpdateResultFields) map[string]any {
	m := map[string]any{
		"identifier":    issue.Identifier,
		"title":         issue.Title,
		"priorityLabel": issue.PriorityLabel,
		"status":        map[string]any{"id": issue.State.Id, "name": issue.State.Name, "type": issue.State.Type},
	}
	if issue.Assignee != nil {
		m["assignee"] = map[string]any{"id": issue.Assignee.Id, "name": issue.Assignee.Name}
	}
	if issue.Project != nil {
		m["project"] = map[string]any{"id": issue.Project.Id, "name": issue.Project.Name}
	}
	if issue.Cycle != nil {
		cycle := map[string]any{"id": issue.Cycle.Id, "number": issue.Cycle.Number}
		setIfNotNil(cycle, "name", issue.Cycle.Name)
		m["cycle"] = cycle
	}
	if issue.Parent != nil {
		m["parent"] = map[string]any{"id": issue.Parent.Id, "identifier": issue.Parent.Identifier}
	}
	if len(issue.Labels.Nodes) > 0 {
		labels := make([]map[string]any, len(issue.Labels.Nodes))
		for i, l := range issue.Labels.Nodes {
			labels[i] = map[string]any{"id": l.Id, "name": l.Name}
		}
		m["labels"] = labels
	}
	setIfNotNil(m, "estimate", issue.Estimate)
	setIfNotNil(m, "dueDate", issue.DueDate)
	return m
}
