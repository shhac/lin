package mappers

import "github.com/shhac/lin/internal/linear"

// HistoryNode is the genqlient type for a single issue history entry.
type HistoryNode = linear.IssueHistoryIssueHistoryIssueHistoryConnectionNodesIssueHistory

// MapHistoryEntry converts a single issue history node into the standard output map.
func MapHistoryEntry(h HistoryNode) map[string]any {
	m := map[string]any{
		"id":                 h.Id,
		"createdAt":          h.CreatedAt,
		"fromPriority":       h.FromPriority,
		"toPriority":         h.ToPriority,
		"fromEstimate":       h.FromEstimate,
		"toEstimate":         h.ToEstimate,
		"fromTitle":          h.FromTitle,
		"toTitle":            h.ToTitle,
		"fromDueDate":        h.FromDueDate,
		"toDueDate":          h.ToDueDate,
		"updatedDescription": h.UpdatedDescription,
		"archived":           h.Archived,
		"trashed":            h.Trashed,
		"autoArchived":       h.AutoArchived,
		"autoClosed":         h.AutoClosed,
	}

	if h.Actor != nil {
		m["actor"] = map[string]any{"id": h.Actor.Id, "name": h.Actor.Name}
	}
	if h.FromState != nil {
		m["fromState"] = map[string]any{"id": h.FromState.Id, "name": h.FromState.Name}
	}
	if h.ToState != nil {
		m["toState"] = map[string]any{"id": h.ToState.Id, "name": h.ToState.Name}
	}
	if h.FromAssignee != nil {
		m["fromAssignee"] = map[string]any{"id": h.FromAssignee.Id, "name": h.FromAssignee.Name}
	}
	if h.ToAssignee != nil {
		m["toAssignee"] = map[string]any{"id": h.ToAssignee.Id, "name": h.ToAssignee.Name}
	}
	if h.FromProject != nil {
		m["fromProject"] = map[string]any{"id": h.FromProject.Id, "name": h.FromProject.Name}
	}
	if h.ToProject != nil {
		m["toProject"] = map[string]any{"id": h.ToProject.Id, "name": h.ToProject.Name}
	}

	addedLabels := make([]map[string]any, len(h.AddedLabels))
	for i, l := range h.AddedLabels {
		addedLabels[i] = map[string]any{"id": l.Id, "name": l.Name}
	}
	m["addedLabels"] = addedLabels

	removedLabels := make([]map[string]any, len(h.RemovedLabels))
	for i, l := range h.RemovedLabels {
		removedLabels[i] = map[string]any{"id": l.Id, "name": l.Name}
	}
	m["removedLabels"] = removedLabels

	return m
}
