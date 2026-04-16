package team

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/estimates"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerGet(team *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Team details and members",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveTeam(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.TeamGet(ctx, client, resolved.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}
			t := resp.Team

			var validValues []int
			var display *string
			if t.IssueEstimationType != "notUsed" {
				cfg := estimates.BuildConfig(t.IssueEstimationType, t.IssueEstimationAllowZero, t.IssueEstimationExtended)
				validValues = estimates.ValidEstimates(cfg)
				d := estimates.FormatScale(t.IssueEstimationType, validValues)
				display = &d
			}

			members := make([]map[string]any, len(t.Members.Nodes))
			for i, m := range t.Members.Nodes {
				members[i] = map[string]any{
					"id":    m.Id,
					"name":  m.Name,
					"email": m.Email,
				}
			}

			output.PrintJSON(map[string]any{
				"id":          t.Id,
				"name":        t.Name,
				"key":         t.Key,
				"description": t.Description,
				"estimates": map[string]any{
					"type":        t.IssueEstimationType,
					"allowZero":   t.IssueEstimationAllowZero,
					"extended":    t.IssueEstimationExtended,
					"default":     t.DefaultIssueEstimate,
					"validValues": validValues,
					"display":     display,
				},
				"members": members,
			})
		},
	}
	team.AddCommand(cmd)
}
