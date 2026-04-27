package issue

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/estimates"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerUpdateEstimate(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "estimate <id> <value>",
		Short: "Update issue estimate (validated against team scale)",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			est, err := strconv.Atoi(args[1])
			if err != nil {
				output.PrintErrorf("Invalid estimate: %q. Must be a number.", args[1])
			}

			client := linear.GetClient()
			ctx := context.Background()

			teamID := resolveIssueTeamID(ctx, client, args[0])
			if teamID == "" {
				output.PrintError("Could not resolve team for this issue.")
			}

			teamDetail, err := linear.TeamGet(ctx, client, teamID)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			cfg := estimates.BuildConfig(
				teamDetail.Team.IssueEstimationType,
				teamDetail.Team.IssueEstimationAllowZero,
				teamDetail.Team.IssueEstimationExtended,
			)
			if validateErr := estimates.Validate(cfg, est); validateErr != nil {
				output.PrintError(validateErr.Error())
			}

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{Estimate: &est})
			if err != nil {
				output.HandleGraphQLError(err)
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}
