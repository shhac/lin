package cycle

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
)

func registerGet(cycle *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>...",
		Short: "Cycle details with issues",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return shared.GetEntities(args, func(client graphql.Client, id string) (any, error) {
				resp, err := linear.CycleGet(context.Background(), client, id)
				if err != nil {
					return nil, apierrors.ClassifyGraphQLError(err)
				}

				c := resp.Cycle
				issues := make([]map[string]any, len(c.Issues.Nodes))
				for i, issue := range c.Issues.Nodes {
					issues[i] = mappers.MapIssueSummary(mappers.FromIssueSummaryFields(issue.IssueSummaryFields))
				}

				result := mapCycleSummary(c.Id, c.Number, c.Name, c.StartsAt, c.EndsAt)
				result["issues"] = issues
				return result, nil
			})
		},
	}
	cycle.AddCommand(cmd)
}
