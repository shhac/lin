package mappers

import "github.com/shhac/lin/internal/linear"

// MapIssueDetail builds the standard output map for an issue get command.
// Comments and attachments are included in the IssueGet query response.
func MapIssueDetail(issue linear.IssueGetIssue) map[string]any {
	statusObj := map[string]any{
		"id":   issue.State.Id,
		"name": issue.State.Name,
		"type": issue.State.Type,
	}

	var assigneeObj map[string]any
	if issue.Assignee != nil {
		assigneeObj = map[string]any{"id": issue.Assignee.Id, "name": issue.Assignee.Name}
	}

	teamObj := map[string]any{"id": issue.Team.Id, "key": issue.Team.Key, "name": issue.Team.Name}

	var projectObj map[string]any
	if issue.Project != nil {
		projectObj = map[string]any{"id": issue.Project.Id, "name": issue.Project.Name}
	}

	labels := make([]map[string]any, len(issue.Labels.Nodes))
	for i, l := range issue.Labels.Nodes {
		labels[i] = map[string]any{"id": l.Id, "name": l.Name}
	}

	attachmentMaps := make([]map[string]any, len(issue.Attachments.Nodes))
	for i, a := range issue.Attachments.Nodes {
		attachmentMaps[i] = map[string]any{
			"title":      a.Title,
			"url":        a.Url,
			"sourceType": a.SourceType,
		}
	}

	var parentObj map[string]any
	if issue.Parent != nil {
		parentObj = map[string]any{"id": issue.Parent.Id, "identifier": issue.Parent.Identifier}
	}

	return map[string]any{
		"id":            issue.Id,
		"identifier":    issue.Identifier,
		"url":           issue.Url,
		"title":         issue.Title,
		"description":   issue.Description,
		"branchName":    issue.BranchName,
		"status":        statusObj,
		"assignee":      assigneeObj,
		"team":          teamObj,
		"project":       projectObj,
		"priority":      issue.Priority,
		"priorityLabel": issue.PriorityLabel,
		"commentCount":  len(issue.Comments.Nodes),
		"labels":        labels,
		"attachments":   attachmentMaps,
		"parent":        parentObj,
		"estimate":      issue.Estimate,
		"dueDate":       issue.DueDate,
		"createdAt":     issue.CreatedAt,
		"updatedAt":     issue.UpdatedAt,
	}
}
