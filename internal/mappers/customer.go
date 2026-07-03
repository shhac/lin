package mappers

import "github.com/shhac/lin/internal/linear"

// MapCustomerSummary converts the CustomerSummaryFields fragment (used by
// CustomerList and customer search) into the standard summary output map.
func MapCustomerSummary(f linear.CustomerSummaryFields) map[string]any {
	out := map[string]any{
		"id":                   f.Id,
		"name":                 f.Name,
		"slugId":               f.SlugId,
		"url":                  f.Url,
		"approximateNeedCount": f.ApproximateNeedCount,
		"status":               map[string]any{"id": f.Status.Id, "name": f.Status.Name},
	}
	if f.Revenue != nil {
		out["revenue"] = *f.Revenue
	}
	if f.Owner != nil {
		out["owner"] = map[string]any{"id": f.Owner.Id, "name": f.Owner.Name}
	}
	if f.Tier != nil {
		out["tier"] = map[string]any{"id": f.Tier.Id, "displayName": f.Tier.DisplayName}
	}
	return out
}

// MapCustomerDetail builds the output map for a customer get command.
func MapCustomerDetail(c linear.CustomerGetCustomer) map[string]any {
	out := map[string]any{
		"id":                   c.Id,
		"name":                 c.Name,
		"slugId":               c.SlugId,
		"url":                  c.Url,
		"approximateNeedCount": c.ApproximateNeedCount,
		"status":               map[string]any{"id": c.Status.Id, "name": c.Status.Name},
		"domains":              c.Domains,
		"externalIds":          c.ExternalIds,
		"createdAt":            c.CreatedAt,
		"updatedAt":            c.UpdatedAt,
	}
	if c.Revenue != nil {
		out["revenue"] = *c.Revenue
	}
	if c.Size != nil {
		out["size"] = *c.Size
	}
	if c.Owner != nil {
		out["owner"] = map[string]any{"id": c.Owner.Id, "name": c.Owner.Name}
	}
	if c.Tier != nil {
		out["tier"] = map[string]any{"id": c.Tier.Id, "displayName": c.Tier.DisplayName}
	}
	return out
}

// MapCustomerNeedSummary converts the CustomerNeedSummaryFields fragment (used
// by CustomerNeeds, IssueNeeds, and ProjectNeeds) into the standard summary
// output map. A need links to either an issue or a project, never both, so the
// mapper emits whichever is set.
func MapCustomerNeedSummary(f linear.CustomerNeedSummaryFields) map[string]any {
	out := map[string]any{
		"id":        f.Id,
		"important": f.Priority == 1,
		"createdAt": f.CreatedAt,
	}
	if f.Body != nil {
		out["body"] = *f.Body
	}
	if f.Url != nil {
		out["url"] = *f.Url
	}
	if f.Customer != nil {
		out["customer"] = map[string]any{"id": f.Customer.Id, "name": f.Customer.Name}
	}
	if f.Issue != nil {
		out["issue"] = map[string]any{"identifier": f.Issue.Identifier, "title": f.Issue.Title}
	}
	if f.Project != nil {
		out["project"] = map[string]any{"id": f.Project.Id, "name": f.Project.Name}
	}
	return out
}
