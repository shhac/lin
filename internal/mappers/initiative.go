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

// MapInitiativeDetail builds the output map for an initiative get command.
func MapInitiativeDetail(i linear.InitiativeGetInitiative) map[string]any {
	var owner any
	if i.Owner != nil {
		owner = map[string]any{
			"id":   i.Owner.Id,
			"name": i.Owner.Name,
		}
	}

	result := map[string]any{
		"id":     i.Id,
		"slugId": i.SlugId,
		"url":    i.Url,
		"name":   i.Name,
		"status": i.Status,
		"owner":  owner,
		"creator": map[string]any{
			"id":   i.Creator.Id,
			"name": i.Creator.Name,
		},
		"createdAt": i.CreatedAt,
		"updatedAt": i.UpdatedAt,
	}
	if i.Description != nil {
		result["description"] = *i.Description
	}
	if i.Content != nil {
		result["content"] = *i.Content
	}
	if i.Health != nil {
		result["health"] = *i.Health
	}
	if i.TargetDate != nil {
		result["targetDate"] = *i.TargetDate
	}
	if i.StartedAt != nil {
		result["startedAt"] = *i.StartedAt
	}
	if i.CompletedAt != nil {
		result["completedAt"] = *i.CompletedAt
	}
	return result
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
	setIfNotNil(m, "health", s.Health)
	setIfNotNil(m, "targetDate", s.TargetDate)
	return m
}
