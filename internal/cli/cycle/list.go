package cycle

import (
	"context"
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

	cmd := &cobra.Command{
		Use:   "list <team>",
		Short: "List cycles",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(_ *cobra.Command, args []string) {
		client := linear.GetClient()
		ctx := context.Background()

		resolved, err := resolvers.ResolveTeam(client, args[0])
		if err != nil {
			output.PrintError(err.Error())
		}

		if current {
			resp, err := linear.TeamActiveCycle(ctx, client, resolved.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}
			c := resp.Team.ActiveCycle
			if c == nil {
				output.PrintJSON([]any{})
				return
			}
			output.PrintJSON([]any{mapCycleSummary(c.Id, c.Number, c.Name, c.StartsAt, c.EndsAt)})
			return
		}

		resp, err := linear.TeamCycles(ctx, client, resolved.ID, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		nodes := resp.Team.Cycles.Nodes
		now := time.Now()

		if next {
			if n, ok := findNextCycle(nodes, now); ok {
				output.PrintJSON([]any{mapCycleSummary(n.Id, n.Number, n.Name, n.StartsAt, n.EndsAt)})
			} else {
				output.PrintJSON([]any{})
			}
			return
		}

		if previous {
			if p, ok := findPreviousCycle(nodes, now); ok {
				output.PrintJSON([]any{mapCycleSummary(p.Id, p.Number, p.Name, p.StartsAt, p.EndsAt)})
			} else {
				output.PrintJSON([]any{})
			}
			return
		}

		items := make([]map[string]any, len(nodes))
		for i, c := range nodes {
			items[i] = mapCycleSummary(c.Id, c.Number, c.Name, c.StartsAt, c.EndsAt)
		}

		output.PrintPage(items, resp.Team.Cycles.PageInfo.HasNextPage, resp.Team.Cycles.PageInfo.EndCursor)
	}

	cmd.Flags().BoolVar(&current, "current", false, "Show only current cycle")
	cmd.Flags().BoolVar(&next, "next", false, "Show only next cycle")
	cmd.Flags().BoolVar(&previous, "previous", false, "Show only previous cycle")
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

type cycleNode = linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle

// findNextCycle returns the cycle with the earliest start time after now.
func findNextCycle(nodes []cycleNode, now time.Time) (cycleNode, bool) {
	return pickCycle(nodes, func(c cycleNode) string { return c.StartsAt }, func(t time.Time) bool { return t.After(now) }, time.Time.Before)
}

// findPreviousCycle returns the cycle with the latest end time before now.
func findPreviousCycle(nodes []cycleNode, now time.Time) (cycleNode, bool) {
	return pickCycle(nodes, func(c cycleNode) string { return c.EndsAt }, func(t time.Time) bool { return t.Before(now) }, time.Time.After)
}

// pickCycle is a single-pass argmin/argmax over cycles. timeOf extracts the
// time field; keep filters which entries to consider; better reports whether
// a is preferred over b (returns true ⇒ a wins).
func pickCycle(
	nodes []cycleNode,
	timeOf func(cycleNode) string,
	keep func(time.Time) bool,
	better func(a, b time.Time) bool,
) (cycleNode, bool) {
	var (
		best   cycleNode
		bestAt time.Time
		found  bool
	)
	for _, c := range nodes {
		t, err := time.Parse(time.RFC3339, timeOf(c))
		if err != nil || !keep(t) {
			continue
		}
		if !found || better(t, bestAt) {
			best, bestAt, found = c, t, true
		}
	}
	return best, found
}
