package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerGet(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get issue details: title, description, status, assignee, labels, relationships",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()
			id := args[0]

			resp, err := linear.IssueGet(ctx, client, id)
			if err != nil {
				output.PrintError(err.Error())
			}

			commentsResp, err := linear.IssueComments(ctx, client, id, 250, nil)
			if err != nil {
				output.PrintError(err.Error())
			}

			attachResp, err := linear.IssueAttachments(ctx, client, id)
			if err != nil {
				output.PrintError(err.Error())
			}

			i := resp.Issue

			statusObj := map[string]any{
				"id":   i.State.Id,
				"name": i.State.Name,
				"type": i.State.Type,
			}

			var assigneeObj map[string]any
			if i.Assignee != nil {
				assigneeObj = map[string]any{"id": i.Assignee.Id, "name": i.Assignee.Name}
			}

			teamObj := map[string]any{"id": i.Team.Id, "key": i.Team.Key, "name": i.Team.Name}

			var projectObj map[string]any
			if i.Project != nil {
				projectObj = map[string]any{"id": i.Project.Id, "name": i.Project.Name}
			}

			labels := make([]map[string]any, len(i.Labels.Nodes))
			for j, l := range i.Labels.Nodes {
				labels[j] = map[string]any{"id": l.Id, "name": l.Name}
			}

			attachments := make([]map[string]any, len(attachResp.Issue.Attachments.Nodes))
			for j, a := range attachResp.Issue.Attachments.Nodes {
				attachments[j] = map[string]any{
					"title":      a.Title,
					"url":        a.Url,
					"sourceType": a.SourceType,
				}
			}

			var parentObj map[string]any
			if i.Parent != nil {
				parentObj = map[string]any{"id": i.Parent.Id, "identifier": i.Parent.Identifier}
			}

			result := map[string]any{
				"id":            i.Id,
				"identifier":    i.Identifier,
				"url":           i.Url,
				"title":         i.Title,
				"description":   i.Description,
				"branchName":    i.BranchName,
				"status":        statusObj,
				"assignee":      assigneeObj,
				"team":          teamObj,
				"project":       projectObj,
				"priority":      i.Priority,
				"priorityLabel": i.PriorityLabel,
				"commentCount":  len(commentsResp.Issue.Comments.Nodes),
				"labels":        labels,
				"attachments":   attachments,
				"parent":        parentObj,
				"estimate":      i.Estimate,
				"dueDate":       i.DueDate,
				"createdAt":     i.CreatedAt,
				"updatedAt":     i.UpdatedAt,
			}

			output.PrintJSON(result)
		},
	}

	parent.AddCommand(cmd)
}
