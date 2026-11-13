package issue

import (
	"context"
	"strconv"

	"github.com/Khan/genqlient/graphql"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/estimates"
	"github.com/shhac/lin/internal/linear"
)

func estimateInput(ctx context.Context, client graphql.Client, issueID, value string) (linear.IssueUpdateInput, error) {
	est, err := strconv.Atoi(value)
	if err != nil {
		return linear.IssueUpdateInput{}, apierrors.Newf(apierrors.FixableByAgent, "invalid estimate: %q, must be a number", value)
	}
	teamID, err := resolveIssueTeamID(ctx, client, issueID)
	if err != nil {
		return linear.IssueUpdateInput{}, err
	}
	teamDetail, err := linear.TeamGet(ctx, client, teamID)
	if err != nil {
		return linear.IssueUpdateInput{}, err
	}
	cfg := estimates.Config{
		Type:      teamDetail.Team.IssueEstimationType,
		AllowZero: teamDetail.Team.IssueEstimationAllowZero,
		Extended:  teamDetail.Team.IssueEstimationExtended,
	}
	if err := estimates.Validate(cfg, est); err != nil {
		return linear.IssueUpdateInput{}, apierrors.Wrap(err, apierrors.FixableByAgent)
	}
	return linear.IssueUpdateInput{Estimate: &est}, nil
}
