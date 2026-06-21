package label

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerGet(label *cobra.Command) {
	var (
		typeFlag string
		teamFlag string
	)

	cmd := &cobra.Command{
		Use:   "get <id|name>...",
		Short: "Get a single label by ID or exact name",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := validateLabelTypeFlags(typeFlag, teamFlag); err != nil {
				return err
			}

			// Team resolution is command-level: a bad --team flag aborts the
			// whole batch rather than producing per-item @unresolved records.
			teamID, err := resolvers.ResolveOptionalTeamID(linear.GetClient(), teamFlag)
			if err != nil {
				return err
			}

			return shared.GetEntities(args, func(client graphql.Client, input string) (any, error) {
				if typeFlag == labelTypeProject {
					return getProjectLabel(context.Background(), client, input)
				}
				return getIssueLabel(context.Background(), client, input, teamID)
			})
		},
	}

	addTypeFlag(cmd, &typeFlag)
	cmd.Flags().StringVar(&teamFlag, "team", "", "Disambiguate by team (issue labels only; key, name, or UUID)")
	label.AddCommand(cmd)
}

func getIssueLabel(ctx context.Context, client graphql.Client, input, teamID string) (any, error) {
	var filter *linear.IssueLabelFilter
	if filters.IsUUID(input) {
		filter = &linear.IssueLabelFilter{Id: &linear.IDComparator{Eq: ptr.To(input)}}
	} else {
		filter = filters.BuildIssueLabelFilter(filters.IssueLabelFilterOpts{Name: input}, teamID)
	}

	resp, err := linear.IssueLabelList(ctx, client, 50, nil, filter)
	if err != nil {
		return nil, apierrors.ClassifyGraphQLError(err)
	}

	nodes := resp.IssueLabels.Nodes
	if len(nodes) == 0 {
		return nil, labelNotFoundErr("label", input, "try `lin label search` to find candidates")
	}
	if len(nodes) > 1 {
		return nil, ambiguousLabelErr("label", input, len(nodes),
			"pass a UUID or use --team to disambiguate",
			"`lin label list --name <name>` shows all matches with team info")
	}

	return mapIssueLabel(nodes[0].IssueLabelFields), nil
}

func getProjectLabel(ctx context.Context, client graphql.Client, input string) (any, error) {
	var filter *linear.ProjectLabelFilter
	if filters.IsUUID(input) {
		filter = &linear.ProjectLabelFilter{Id: &linear.IDComparator{Eq: ptr.To(input)}}
	} else {
		filter = filters.BuildProjectLabelFilter(filters.ProjectLabelFilterOpts{Name: input})
	}

	resp, err := linear.ProjectLabelList(ctx, client, 50, nil, filter)
	if err != nil {
		return nil, apierrors.ClassifyGraphQLError(err)
	}

	nodes := resp.ProjectLabels.Nodes
	if len(nodes) == 0 {
		return nil, labelNotFoundErr("project label", input,
			"try `lin label search --type project` to find candidates")
	}
	if len(nodes) > 1 {
		return nil, ambiguousLabelErr("project label", input, len(nodes),
			"pass a UUID to disambiguate",
			"`lin label list --type project --name <name>` shows all matches")
	}

	return mapProjectLabel(nodes[0].ProjectLabelFields), nil
}
