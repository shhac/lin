package team

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerStates(team *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "states <team>",
		Short: "List workflow states for a team (discover valid status values)",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			teamInput := args[0]
			client := linear.GetClient()

			team, err := resolvers.ResolveTeam(client, teamInput)
			if err != nil {
				output.PrintError(err.Error())
			}

			filter := &linear.WorkflowStateFilter{
				Team: &linear.TeamFilter{Id: &linear.IDComparator{Eq: ptr.To(team.ID)}},
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
