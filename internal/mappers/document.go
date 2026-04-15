package mappers

// DocSummaryInput holds the fields needed to produce a document summary.
type DocSummaryInput struct {
	ID          string
	SlugId      string
	Title       string
	URL         string
	UpdatedAt   string
	CreatorID   string // empty if no creator
	CreatorName string
	ProjectID   string // empty if no project
	ProjectName string
}

// MapDocSummary converts a DocSummaryInput to the standard output map.
func MapDocSummary(s DocSummaryInput) map[string]any {
	m := map[string]any{
		"id":        s.ID,
		"slugId":    s.SlugId,
		"title":     s.Title,
		"url":       s.URL,
		"updatedAt": s.UpdatedAt,
	}
	if s.CreatorID != "" {
		m["creator"] = map[string]any{"id": s.CreatorID, "name": s.CreatorName}
	}
	if s.ProjectID != "" {
		m["project"] = map[string]any{"id": s.ProjectID, "name": s.ProjectName}
	}
	return m
}
