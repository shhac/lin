package label

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerList(label *cobra.Command) {
	var (
		teamFlag  string
		nameFlag  string
		groupFlag bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels (optionally filtered)",
		Args:  cobra.NoArgs,
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(cmd *cobra.Command, _ []string) {
		client := linear.GetClient()
		ctx := context.Background()

		teamID, err := resolvers.ResolveOptionalTeamID(client, teamFlag)
		if err != nil {
			output.PrintError(err.Error())
		}

		opts := filters.IssueLabelFilterOpts{Name: nameFlag}
		if cmd.Flags().Changed("is-group") {
			opts.IsGroup = ptr.To(groupFlag)
		}
		filter := filters.BuildIssueLabelFilter(opts, teamID)

		resp, err := linear.IssueLabelList(ctx, client, page.Size(), page.Cursor(), filter)
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]map[string]any, len(resp.IssueLabels.Nodes))
		for i, n := range resp.IssueLabels.Nodes {
			items[i] = mapIssueLabel(n.IssueLabelFields)
		}

		output.PrintPage(items, resp.IssueLabels.PageInfo.HasNextPage, resp.IssueLabels.PageInfo.EndCursor)
	}

	cmd.Flags().StringVar(&teamFlag, "team", "", "Filter by team key, name, or UUID")
	cmd.Flags().StringVar(&nameFlag, "name", "", "Filter by exact label name (case-insensitive)")
	cmd.Flags().BoolVar(&groupFlag, "is-group", false, "Filter to only group labels (--is-group=false for non-groups)")
	label.AddCommand(cmd)
}

func mapIssueLabel(l linear.IssueLabelFields) map[string]any {
	m := map[string]any{
		"id":    l.Id,
		"name":  l.Name,
		"color": l.Color,
	}
	if l.Description != nil && *l.Description != "" {
		m["description"] = *l.Description
	}
	if l.IsGroup {
		m["isGroup"] = true
	}
	if l.Team != nil {
		m["team"] = map[string]any{
			"id":   l.Team.Id,
			"key":  l.Team.Key,
			"name": l.Team.Name,
		}
	}
	if l.Parent != nil {
		m["parent"] = map[string]any{
			"id":   l.Parent.Id,
			"name": l.Parent.Name,
		}
	}
	return m
}
