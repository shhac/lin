package label

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerGet(label *cobra.Command) {
	var (
		typeFlag string
		teamFlag string
	)

	cmd := &cobra.Command{
		Use:   "get <id|name>",
		Short: "Get a single label by ID or exact name",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if err := validateType(typeFlag); err != nil {
				output.WriteError(err)
			}
			if err := rejectTeamForProject(typeFlag, teamFlag); err != nil {
				output.WriteError(err)
			}

			client := linear.GetClient()
			ctx := context.Background()
			input := args[0]

			if typeFlag == labelTypeProject {
				getProjectLabel(ctx, client, input)
				return
			}

			teamID, err := resolvers.ResolveOptionalTeamID(client, teamFlag)
			if err != nil {
				output.PrintError(err.Error())
			}
			getIssueLabel(ctx, client, input, teamID)
		},
	}

	addTypeFlag(cmd, &typeFlag)
	cmd.Flags().StringVar(&teamFlag, "team", "", "Disambiguate by team (issue labels only; key, name, or UUID)")
	label.AddCommand(cmd)
}

func getIssueLabel(ctx context.Context, client graphql.Client, input, teamID string) {
	var filter *linear.IssueLabelFilter
	if filters.IsUUID(input) {
		filter = &linear.IssueLabelFilter{Id: &linear.IDComparator{Eq: ptr.To(input)}}
	} else {
		filter = filters.BuildIssueLabelFilter(filters.IssueLabelFilterOpts{Name: input}, teamID)
	}

	resp, err := linear.IssueLabelList(ctx, client, 50, nil, filter)
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

	output.PrintJSON(mapIssueLabel(nodes[0].IssueLabelFields))
}

func getProjectLabel(ctx context.Context, client graphql.Client, input string) {
	var filter *linear.ProjectLabelFilter
	if filters.IsUUID(input) {
		filter = &linear.ProjectLabelFilter{Id: &linear.IDComparator{Eq: ptr.To(input)}}
	} else {
		filter = filters.BuildProjectLabelFilter(filters.ProjectLabelFilterOpts{Name: input})
	}

	resp, err := linear.ProjectLabelList(ctx, client, 50, nil, filter)
	if err != nil {
		output.HandleGraphQLError(err)
	}

	nodes := resp.ProjectLabels.Nodes
	if len(nodes) == 0 {
		output.WriteError(apierrors.Newf(apierrors.FixableByAgent, "project label not found: %q", input).
			WithHint("try `lin label search --type project` to find candidates"))
	}
	if len(nodes) > 1 {
		output.WriteError(apierrors.Newf(apierrors.FixableByAgent,
			"%d project labels match %q — pass a UUID to disambiguate", len(nodes), input).
			WithHint("`lin label list --type project --name <name>` shows all matches"))
	}

	output.PrintJSON(mapProjectLabel(nodes[0].ProjectLabelFields))
}
