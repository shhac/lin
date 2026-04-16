package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
)

func registerHistory(parent *cobra.Command) {
	var (
		limit  string
		cursor string
	)

	cmd := &cobra.Command{
		Use:   "history <issue-id>",
		Short: "List activity history for an issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			pageSize := output.ResolvePageSize(limit)
			afterPtr := output.ResolveCursor(cursor)

			resp, err := linear.IssueHistory(ctx, client, args[0], pageSize, afterPtr)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.Issue.History.Nodes))
			for i, h := range resp.Issue.History.Nodes {
				m := map[string]any{
					"id":                 h.Id,
					"createdAt":          h.CreatedAt,
					"fromPriority":       h.FromPriority,
					"toPriority":         h.ToPriority,
					"fromEstimate":       h.FromEstimate,
					"toEstimate":         h.ToEstimate,
					"fromTitle":          h.FromTitle,
					"toTitle":            h.ToTitle,
					"fromDueDate":        h.FromDueDate,
					"toDueDate":          h.ToDueDate,
					"updatedDescription": h.UpdatedDescription,
					"archived":           h.Archived,
					"trashed":            h.Trashed,
					"autoArchived":       h.AutoArchived,
					"autoClosed":         h.AutoClosed,
				}

				if h.Actor != nil {
					m["actor"] = map[string]any{"id": h.Actor.Id, "name": h.Actor.Name}
				}
				if h.FromState != nil {
					m["fromState"] = map[string]any{"id": h.FromState.Id, "name": h.FromState.Name}
				}
				if h.ToState != nil {
					m["toState"] = map[string]any{"id": h.ToState.Id, "name": h.ToState.Name}
				}
				if h.FromAssignee != nil {
					m["fromAssignee"] = map[string]any{"id": h.FromAssignee.Id, "name": h.FromAssignee.Name}
				}
				if h.ToAssignee != nil {
					m["toAssignee"] = map[string]any{"id": h.ToAssignee.Id, "name": h.ToAssignee.Name}
				}
				if h.FromProject != nil {
					m["fromProject"] = map[string]any{"id": h.FromProject.Id, "name": h.FromProject.Name}
				}
				if h.ToProject != nil {
					m["toProject"] = map[string]any{"id": h.ToProject.Id, "name": h.ToProject.Name}
				}

				addedLabels := make([]map[string]any, len(h.AddedLabels))
				for j, l := range h.AddedLabels {
					addedLabels[j] = map[string]any{"id": l.Id, "name": l.Name}
				}
				m["addedLabels"] = addedLabels

				removedLabels := make([]map[string]any, len(h.RemovedLabels))
				for j, l := range h.RemovedLabels {
					removedLabels[j] = map[string]any{"id": l.Id, "name": l.Name}
				}
				m["removedLabels"] = removedLabels

				items[i] = m
			}

			pi := resp.Issue.History.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}
