package label

import (
	"context"

	"github.com/spf13/cobra"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerGet(label *cobra.Command) {
	var teamFlag string

	cmd := &cobra.Command{
		Use:   "get <id|name>",
		Short: "Get a single label by ID or exact name",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()
			input := args[0]

			var teamID string
			if teamFlag != "" {
				resolved, err := resolvers.ResolveTeam(client, teamFlag)
				if err != nil {
					output.PrintError(err.Error())
				}
				teamID = resolved.ID
			}

			var filter *linear.IssueLabelFilter
			if filters.IsUUID(input) {
				filter = &linear.IssueLabelFilter{Id: &linear.IDComparator{Eq: ptr.To(input)}}
			} else {
				filter = filters.BuildIssueLabelFilter(filters.LabelFilterOpts{Name: input}, teamID)
			}

			resp, err := linear.LabelList(ctx, client, 50, nil, filter)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			nodes := resp.IssueLabels.Nodes
			if len(nodes) == 0 {
				output.WriteError(apierrors.Newf(apierrors.FixableByAgent, "label not found: %q", input).
					WithHint("try `lin label search` to find candidates"))
			}
			if len(nodes) > 1 {
				output.WriteError(apierrors.Newf(apierrors.FixableByAgent,
					"%d labels match %q — pass a UUID or use --team to disambiguate", len(nodes), input).
					WithHint("`lin label list --name <name>` shows all matches with team info"))
			}

			output.PrintJSON(mapLabel(nodes[0].LabelFields))
		},
	}

	cmd.Flags().StringVar(&teamFlag, "team", "", "Disambiguate by team (key, name, or UUID)")
	label.AddCommand(cmd)
}
