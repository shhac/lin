package team

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/estimates"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/resolvers"
)

func registerGet(team *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>...",
		Short: "Team details and members",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return shared.GetEntities(args, func(client graphql.Client, id string) (any, error) {
				resolved, err := resolvers.ResolveTeam(client, id)
				if err != nil {
					return nil, err
				}

				resp, err := linear.TeamGet(context.Background(), client, resolved.ID)
				if err != nil {
					return nil, apierrors.ClassifyGraphQLError(err)
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

				return map[string]any{
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
				}, nil
			})
		},
	}
	team.AddCommand(cmd)
}
