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

	registerUpdateStatus(update)
	registerUpdateAssignee(update)
	registerUpdatePriority(update)
	registerUpdateProject(update)
	registerUpdateLabels(update)
	registerUpdateEstimate(update)

	registerSimpleIssueUpdate(update, "title <id> <new-title>", "Update issue title",
		func(v string) linear.IssueUpdateInput { return linear.IssueUpdateInput{Title: &v} })
	registerSimpleIssueUpdate(update, "description <id> <description>", "Update issue description",
		func(v string) linear.IssueUpdateInput { return linear.IssueUpdateInput{Description: &v} })
	registerSimpleIssueUpdate(update, "due-date <id> <date>", "Update issue due date",
		func(v string) linear.IssueUpdateInput { return linear.IssueUpdateInput{DueDate: &v} })
	registerSimpleIssueUpdate(update, "cycle <id> <cycle-id>", "Move issue to a cycle",
		func(v string) linear.IssueUpdateInput { return linear.IssueUpdateInput{CycleId: &v} })
	registerSimpleIssueUpdate(update, "parent <id> <parent-id>", "Set parent issue (make sub-issue)",
		func(v string) linear.IssueUpdateInput { return linear.IssueUpdateInput{ParentId: &v} })
}

func registerSimpleIssueUpdate(parent *cobra.Command, use, short string, buildInput func(string) linear.IssueUpdateInput) {
	parent.AddCommand(&cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueUpdate(ctx, client, args[0], buildInput(args[1]))
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
