package mappers

import "github.com/shhac/lin/internal/linear"

// InitiativeSummaryInput holds the fields needed to produce an initiative summary.
// Callers populate this from their specific genqlient response types.
type InitiativeSummaryInput struct {
	ID         string
	SlugId     string
	URL        string
	Name       string
	Status     linear.InitiativeStatus
	Health     *linear.InitiativeUpdateHealthType
	TargetDate *string
	OwnerName  *string
}

// FromInitiativeListFields converts a genqlient InitiativeList node into an InitiativeSummaryInput.
func FromInitiativeListFields(n linear.InitiativeListInitiativesInitiativeConnectionNodesInitiative) InitiativeSummaryInput {
	input := InitiativeSummaryInput{
		ID:         n.Id,
		SlugId:     n.SlugId,
		URL:        n.Url,
		Name:       n.Name,
		Status:     n.Status,
		Health:     n.Health,
		TargetDate: n.TargetDate,
	}
	if n.Owner != nil {
		input.OwnerName = &n.Owner.Name
	}
	return input
}

// MapInitiativeSummary converts an InitiativeSummaryInput to the standard output map.
func MapInitiativeSummary(s InitiativeSummaryInput) map[string]any {
	m := map[string]any{
		"id":     s.ID,
		"slugId": s.SlugId,
		"url":    s.URL,
		"name":   s.Name,
		"status": s.Status,
		"owner":  s.OwnerName,
	}
	if s.Health != nil {
		m["health"] = *s.Health
	}
	if s.TargetDate != nil {
		m["targetDate"] = *s.TargetDate
	}
	return m
}
