package team

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerStates(team *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "states <team>",
		Short: "List workflow states for a team (discover valid status values)",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			teamInput := args[0]
			client := linear.GetClient()

			filter := &linear.WorkflowStateFilter{
				Team: filters.BuildTeamFilter(teamInput),
			}
			resp, err := linear.WorkflowStates(context.Background(), client, filter)
			if err != nil {
				output.PrintError(err.Error())
			}

			nodes := resp.WorkflowStates.Nodes
			if len(nodes) == 0 {
				output.PrintError(fmt.Sprintf("No workflow states found for team %q.", teamInput))
			}

			sort.Slice(nodes, func(i, j int) bool {
				return nodes[i].Position < nodes[j].Position
			})

			items := make([]map[string]any, len(nodes))
			for i, s := range nodes {
				items[i] = map[string]any{
					"id":       s.Id,
					"name":     s.Name,
					"type":     s.Type,
					"color":    s.Color,
					"position": s.Position,
				}
			}

			output.PrintJSON(items)
		},
	}
	team.AddCommand(cmd)
}
