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
		groupSet  bool
		limit     string
		cursor    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels (optionally filtered)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			client := linear.GetClient()
			ctx := context.Background()
			pageSize := output.ResolvePageSize(limit)
			after := output.ResolveCursor(cursor)

			groupSet = cmd.Flags().Changed("is-group")

			var teamID string
			if teamFlag != "" {
				resolved, err := resolvers.ResolveTeam(client, teamFlag)
				if err != nil {
					output.PrintError(err.Error())
				}
				teamID = resolved.ID
			}

			opts := filters.LabelFilterOpts{Name: nameFlag}
			if groupSet {
				opts.IsGroup = ptr.To(groupFlag)
			}
			filter := filters.BuildIssueLabelFilter(opts, teamID)

			resp, err := linear.LabelList(ctx, client, pageSize, after, filter)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]map[string]any, len(resp.IssueLabels.Nodes))
			for i, n := range resp.IssueLabels.Nodes {
				items[i] = mapLabel(n.LabelFields)
			}

			pi := resp.IssueLabels.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&teamFlag, "team", "", "Filter by team key, name, or UUID")
	cmd.Flags().StringVar(&nameFlag, "name", "", "Filter by exact label name (case-insensitive)")
	cmd.Flags().BoolVar(&groupFlag, "is-group", false, "Filter to only group labels (--is-group=false for non-groups)")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	label.AddCommand(cmd)
}

func mapLabel(l linear.LabelFields) map[string]any {
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
