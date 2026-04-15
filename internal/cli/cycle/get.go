package cycle

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerGet(cycle *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Cycle details with issues",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()

			resp, err := linear.CycleGet(context.Background(), client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			c := resp.Cycle
			issues := make([]map[string]any, len(c.Issues.Nodes))
			for i, issue := range c.Issues.Nodes {
				f := issue.IssueSummaryFields
				entry := map[string]any{
					"id":            f.Id,
					"identifier":    f.Identifier,
					"title":         f.Title,
					"branchName":    f.BranchName,
					"status":        f.State.Name,
					"statusType":    f.State.Type,
					"priority":      f.Priority,
					"priorityLabel": f.PriorityLabel,
					"team":          f.Team.Key,
				}
				if f.Assignee != nil {
					entry["assignee"] = f.Assignee.Name
					entry["assigneeId"] = f.Assignee.Id
				}
				issues[i] = entry
			}

			result := mapCycleSummary(c.Id, c.Number, c.Name, c.StartsAt, c.EndsAt)
			result["issues"] = issues
			output.PrintJSON(result)
		},
	}
	cycle.AddCommand(cmd)
}
