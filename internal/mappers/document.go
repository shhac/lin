package mappers

import "github.com/shhac/lin/internal/linear"

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

// FromDocSummaryFields converts the genqlient DocSummaryFields fragment
// (used by DocumentList) into a DocSummaryInput.
func FromDocSummaryFields(f linear.DocSummaryFields) DocSummaryInput {
	input := DocSummaryInput{
		ID:        f.Id,
		SlugId:    f.SlugId,
		Title:     f.Title,
		URL:       f.Url,
		UpdatedAt: f.UpdatedAt,
	}
	if f.Creator != nil {
		input.CreatorID = f.Creator.Id
		input.CreatorName = f.Creator.Name
	}
	if f.Project != nil {
		input.ProjectID = f.Project.Id
		input.ProjectName = f.Project.Name
	}
	return input
}

// FromDocSearchSummaryFields converts the genqlient DocSearchSummaryFields
// fragment (used by DocumentSearch) into a DocSummaryInput.
func FromDocSearchSummaryFields(f linear.DocSearchSummaryFields) DocSummaryInput {
	input := DocSummaryInput{
		ID:        f.Id,
		SlugId:    f.SlugId,
		Title:     f.Title,
		URL:       f.Url,
		UpdatedAt: f.UpdatedAt,
	}
	if f.Creator != nil {
		input.CreatorID = f.Creator.Id
		input.CreatorName = f.Creator.Name
	}
	if f.Project != nil {
		input.ProjectID = f.Project.Id
		input.ProjectName = f.Project.Name
	}
	return input
}
