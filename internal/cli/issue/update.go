package issue

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
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

	output.HandleUnknownCommand(update, "Run 'lin issue usage' for available update subcommands")
}

func resolveIssueTeamID(ctx context.Context, client graphql.Client, issueID string) string {
	resp, err := linear.IssueTeam(ctx, client, issueID)
	if err != nil {
		output.HandleGraphQLError(err)
	}
	return resp.Issue.Team.Id
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
				output.HandleGraphQLError(err)
			}
			output.PrintJSON(map[string]any{"updated": resp.IssueUpdate.Success})
		},
	})
}
