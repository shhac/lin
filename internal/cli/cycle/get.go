package cycle

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
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
				output.HandleGraphQLError(err)
			}

			c := resp.Cycle
			issues := make([]map[string]any, len(c.Issues.Nodes))
			for i, issue := range c.Issues.Nodes {
				f := issue.IssueSummaryFields
				input := mappers.IssueSummaryInput{
					ID:            f.Id,
					Identifier:    f.Identifier,
					Title:         f.Title,
					BranchName:    f.BranchName,
					Priority:      f.Priority,
					PriorityLabel: f.PriorityLabel,
					StateName:     f.State.Name,
					StateType:     f.State.Type,
					TeamKey:       f.Team.Key,
				}
				if f.Assignee != nil {
					input.AssigneeID = f.Assignee.Id
					input.AssigneeName = f.Assignee.Name
				}
				issues[i] = mappers.MapIssueSummary(input)
			}

			result := mapCycleSummary(c.Id, c.Number, c.Name, c.StartsAt, c.EndsAt)
			result["issues"] = issues
			output.PrintJSON(result)
		},
	}
	cycle.AddCommand(cmd)
}
