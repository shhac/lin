package mappers

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
