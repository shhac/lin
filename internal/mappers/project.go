package mappers

import (
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/ptr"
)

// ProjectSummaryInput holds the fields needed to produce a project summary.
type ProjectSummaryInput struct {
	ID         string
	SlugId     string
	URL        string
	Name       string
	State      string
	Progress   float64
	LeadName   string // empty if no lead
	StartDate  string // empty if not set
	TargetDate string // empty if not set
}

// MapProjectSummary converts a ProjectSummaryInput to the standard output map.
func MapProjectSummary(s ProjectSummaryInput) map[string]any {
	m := map[string]any{
		"id":       s.ID,
		"slugId":   s.SlugId,
		"url":      s.URL,
		"name":     s.Name,
		"status":   s.State,
		"progress": s.Progress,
	}
	setIfNotEmpty(m, "lead", s.LeadName)
	setIfNotEmpty(m, "startDate", s.StartDate)
	setIfNotEmpty(m, "targetDate", s.TargetDate)
	return m
}

// MapProjectDetail builds the output map for a project get command.
func MapProjectDetail(p linear.ProjectGetProject) map[string]any {
	result := map[string]any{
		"id":          p.Id,
		"slugId":      p.SlugId,
		"url":         p.Url,
		"name":        p.Name,
		"description": p.Description,
		"content":     p.Content,
		"status":      p.State,
		"progress":    p.Progress,
		"startDate":   p.StartDate,
		"targetDate":  p.TargetDate,
	}

	if p.Lead != nil {
		result["lead"] = map[string]any{
			"id":   p.Lead.Id,
			"name": p.Lead.Name,
		}
	}

	labels := make([]map[string]any, len(p.Labels.Nodes))
	for i, l := range p.Labels.Nodes {
		labels[i] = map[string]any{"id": l.Id, "name": l.Name}
	}
	result["labels"] = labels

	milestones := make([]map[string]any, len(p.ProjectMilestones.Nodes))
	for i, m := range p.ProjectMilestones.Nodes {
		milestones[i] = map[string]any{
			"id":         m.Id,
			"name":       m.Name,
			"targetDate": m.TargetDate,
		}
	}
	result["milestones"] = milestones

	return result
}

// FromProjectSummaryFields converts the genqlient ProjectSummaryFields fragment
// (used by ProjectList, InitiativeProjects) into a ProjectSummaryInput.
func FromProjectSummaryFields(f linear.ProjectSummaryFields) ProjectSummaryInput {
	input := ProjectSummaryInput{
		ID:         f.Id,
		SlugId:     f.SlugId,
		URL:        f.Url,
		Name:       f.Name,
		State:      f.State,
		Progress:   f.Progress,
		StartDate:  ptr.Deref(f.StartDate),
		TargetDate: ptr.Deref(f.TargetDate),
	}
	if f.Lead != nil {
		input.LeadName = f.Lead.Name
	}
	return input
}

// FromProjectUpdateSummary converts the genqlient ProjectUpdateSummaryFields
// fragment (used by ProjectPostList / ProjectPostGet) into the standard output map.
func FromProjectUpdateSummary(f linear.ProjectUpdateSummaryFields) map[string]any {
	m := map[string]any{
		"id":        f.Id,
		"url":       f.Url,
		"health":    string(f.Health),
		"body":      f.Body,
		"createdAt": f.CreatedAt,
		"user":      map[string]any{"id": f.User.Id, "name": f.User.Name},
	}
	setIfNotNil(m, "editedAt", f.EditedAt)
	return m
}

// FromProjectSearchSummaryFields converts the genqlient ProjectSearchSummaryFields
// fragment (used by ProjectSearch) into a ProjectSummaryInput.
func FromProjectSearchSummaryFields(f linear.ProjectSearchSummaryFields) ProjectSummaryInput {
	input := ProjectSummaryInput{
		ID:         f.Id,
		SlugId:     f.SlugId,
		URL:        f.Url,
		Name:       f.Name,
		State:      f.State,
		Progress:   f.Progress,
		StartDate:  ptr.Deref(f.StartDate),
		TargetDate: ptr.Deref(f.TargetDate),
	}
	if f.Lead != nil {
		input.LeadName = f.Lead.Name
	}
	return input
}
