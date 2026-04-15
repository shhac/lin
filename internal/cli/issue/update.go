package issue

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/estimates"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/priorities"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdate(parent *cobra.Command) {
	update := &cobra.Command{
		Use:   "update",
		Short: "Update issue fields",
	}
	parent.AddCommand(update)

	registerUpdateTitle(update)
	registerUpdateStatus(update)
	registerUpdateAssignee(update)
	registerUpdatePriority(update)
	registerUpdateProject(update)
	registerUpdateLabels(update)
	registerUpdateDescription(update)
	registerUpdateDueDate(update)
	registerUpdateCycle(update)
	registerUpdateParent(update)
	registerUpdateEstimate(update)
}

func registerUpdateTitle(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "title <id> <new-title>",
		Short: "Update issue title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()
			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{Title: &args[1]})
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}

func registerUpdateStatus(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "status <id> <new-status>",
		Short: "Update issue status",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			teamResp, err := linear.IssueTeam(ctx, client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}
			teamID := teamResp.Issue.Team.Id
			if teamID == "" {
				output.PrintError("Could not resolve team for this issue.")
			}

			state, err := resolvers.ResolveWorkflowState(client, args[1], teamID)
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{StateId: &state.ID})
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}

func registerUpdateAssignee(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "assignee <id> <user>",
		Short: "Update issue assignee",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			user, err := resolvers.ResolveUser(client, args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{AssigneeId: &user.ID})
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}

func registerUpdatePriority(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "priority <id> <priority>",
		Short: "Update issue priority",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			p, ok := priorities.Resolve(args[1])
			if !ok {
				output.PrintErrorf("Invalid priority: %q. Valid values: %s", args[1], priorities.Values)
			}

			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{Priority: &p})
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}

func registerUpdateProject(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "project <id> <project>",
		Short: "Move issue to project",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			proj, err := resolvers.ResolveProject(client, args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{ProjectId: &proj.ID})
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}

func registerUpdateLabels(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "labels <id> <labels>",
		Short: "Set issue labels",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			teamResp, err := linear.IssueTeam(ctx, client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			labelIds, err := resolvers.ResolveLabels(client, args[1], teamResp.Issue.Team.Id)
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{LabelIds: labelIds})
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}

func registerUpdateDescription(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "description <id> <description>",
		Short: "Update issue description",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{Description: &args[1]})
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}

func registerUpdateDueDate(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "due-date <id> <date>",
		Short: "Update issue due date",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{DueDate: &args[1]})
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}

func registerUpdateCycle(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "cycle <id> <cycle-id>",
		Short: "Move issue to a cycle",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{CycleId: &args[1]})
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}

func registerUpdateParent(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "parent <id> <parent-id>",
		Short: "Set parent issue (make sub-issue)",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueUpdate(ctx, client, args[0], linear.IssueUpdateInput{ParentId: &args[1]})
			if err != nil {
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}

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

			teamResp, err := linear.IssueTeam(ctx, client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}
			teamID := teamResp.Issue.Team.Id
			if teamID == "" {
				output.PrintError("Could not resolve team for this issue.")
			}

			teamDetail, err := linear.TeamGet(ctx, client, teamID)
			if err != nil {
				output.PrintError(err.Error())
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
				output.PrintError(err.Error())
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}
