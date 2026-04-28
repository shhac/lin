package label

import "github.com/shhac/lin/internal/linear"

func mapIssueLabel(l linear.IssueLabelFields) map[string]any {
	m := map[string]any{
		"id":    l.Id,
		"name":  l.Name,
		"color": l.Color,
	}
	if l.Description != nil && *l.Description != "" {
		m["description"] = *l.Description
	}
	if l.IsGroup {
		m["isGroup"] = true
	}
	if l.Team != nil {
		m["team"] = map[string]any{
			"id":   l.Team.Id,
			"key":  l.Team.Key,
			"name": l.Team.Name,
		}
	}
	if l.Parent != nil {
		m["parent"] = map[string]any{
			"id":   l.Parent.Id,
			"name": l.Parent.Name,
		}
	}
	return m
}

func mapProjectLabel(l linear.ProjectLabelFields) map[string]any {
	m := map[string]any{
		"id":    l.Id,
		"name":  l.Name,
		"color": l.Color,
	}
	if l.Description != nil && *l.Description != "" {
		m["description"] = *l.Description
	}
	if l.IsGroup {
		m["isGroup"] = true
	}
	if l.Parent != nil {
		m["parent"] = map[string]any{
			"id":   l.Parent.Id,
			"name": l.Parent.Name,
		}
	}
	return m
}
