package cycle

import (
	"context"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerList(cycle *cobra.Command) {
	var current bool
	var next bool
	var previous bool
	var limit string
	var cursor string

	cmd := &cobra.Command{
		Use:   "list <team>",
		Short: "List cycles",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveTeam(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			if current {
				resp, err := linear.TeamActiveCycle(ctx, client, resolved.ID)
				if err != nil {
					output.PrintError(err.Error())
				}
				c := resp.Team.ActiveCycle
				if c == nil {
					output.PrintJSON([]any{})
					return
				}
				output.PrintJSON([]any{mapCycleSummary(c.Id, c.Number, c.Name, c.StartsAt, c.EndsAt)})
				return
			}

			pageSize := output.ResolvePageSize(limit)
			var after *string
			if cursor != "" {
				after = &cursor
			}

			resp, err := linear.TeamCycles(ctx, client, resolved.ID, pageSize, after)
			if err != nil {
				output.PrintError(err.Error())
			}

			nodes := resp.Team.Cycles.Nodes
			now := time.Now()

			if next {
				type cycleEntry struct {
					node    linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle
					startAt time.Time
				}
				var future []cycleEntry
				for _, c := range nodes {
					t, err := time.Parse(time.RFC3339, c.StartsAt)
					if err != nil {
						continue
					}
					if t.After(now) {
						future = append(future, cycleEntry{node: c, startAt: t})
					}
				}
				if len(future) == 0 {
					output.PrintJSON([]any{})
					return
				}
				sort.Slice(future, func(i, j int) bool {
					return future[i].startAt.Before(future[j].startAt)
				})
				n := future[0].node
				output.PrintJSON([]any{mapCycleSummary(n.Id, n.Number, n.Name, n.StartsAt, n.EndsAt)})
				return
			}

			if previous {
				type cycleEntry struct {
					node  linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle
					endAt time.Time
				}
				var past []cycleEntry
				for _, c := range nodes {
					t, err := time.Parse(time.RFC3339, c.EndsAt)
					if err != nil {
						continue
					}
					if t.Before(now) {
						past = append(past, cycleEntry{node: c, endAt: t})
					}
				}
				if len(past) == 0 {
					output.PrintJSON([]any{})
					return
				}
				sort.Slice(past, func(i, j int) bool {
					return past[i].endAt.After(past[j].endAt)
				})
				p := past[0].node
				output.PrintJSON([]any{mapCycleSummary(p.Id, p.Number, p.Name, p.StartsAt, p.EndsAt)})
				return
			}

			items := make([]map[string]any, len(nodes))
			for i, c := range nodes {
				items[i] = mapCycleSummary(c.Id, c.Number, c.Name, c.StartsAt, c.EndsAt)
			}

			pi := resp.Team.Cycles.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().BoolVar(&current, "current", false, "Show only current cycle")
	cmd.Flags().BoolVar(&next, "next", false, "Show only next cycle")
	cmd.Flags().BoolVar(&previous, "previous", false, "Show only previous cycle")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	cycle.AddCommand(cmd)
}

func mapCycleSummary(id string, number float64, name *string, startsAt, endsAt string) map[string]any {
	return map[string]any{
		"id":       id,
		"number":   number,
		"name":     name,
		"startsAt": startsAt,
		"endsAt":   endsAt,
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
