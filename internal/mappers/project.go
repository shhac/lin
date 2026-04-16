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
	if s.LeadName != "" {
		m["lead"] = s.LeadName
	}
	if s.StartDate != "" {
		m["startDate"] = s.StartDate
	}
	if s.TargetDate != "" {
		m["targetDate"] = s.TargetDate
	}
	return m
}

// FromProjectSummaryFields converts the genqlient ProjectSummaryFields fragment
// (used by ProjectList, RoadmapProjects) into a ProjectSummaryInput.
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
