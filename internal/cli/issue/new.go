package issue

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/priorities"
	"github.com/shhac/lin/internal/resolvers"
)

func registerNew(parent *cobra.Command) {
	var (
		teamFlag     string
		project      string
		assignee     string
		priorityFlag string
		status       string
		labels       string
		description  string
		cycle        string
		parentIssue  string
		estimate     string
	)

	cmd := &cobra.Command{
		Use:   "new <title>",
		Short: "Create a new issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()
			title := args[0]

			var priorityPtr *int
			if priorityFlag != "" {
				p, ok := priorities.Resolve(priorityFlag)
				if !ok {
					output.PrintErrorf("Invalid priority: %q. Valid values: %s", priorityFlag, priorities.Values)
				}
				priorityPtr = &p
			}

			team, err := resolvers.ResolveTeam(client, teamFlag)
			if err != nil {
				output.PrintError(err.Error())
			}

			input := linear.IssueCreateInput{
				Title:    &title,
				TeamId:   team.ID,
				Priority: priorityPtr,
			}

			if status != "" {
				state, err := resolvers.ResolveWorkflowState(client, status, team.ID)
				if err != nil {
					output.PrintError(err.Error())
				}
				input.StateId = &state.ID
			}

			if assignee != "" {
				user, err := resolvers.ResolveUser(client, assignee)
				if err != nil {
					output.PrintError(err.Error())
				}
				input.AssigneeId = &user.ID
			}

			if project != "" {
				proj, err := resolvers.ResolveProject(client, project)
				if err != nil {
					output.PrintError(err.Error())
				}
				input.ProjectId = &proj.ID
			}

			if labels != "" {
				labelIds, err := resolvers.ResolveLabels(client, labels, team.ID)
				if err != nil {
					output.PrintError(err.Error())
				}
				input.LabelIds = labelIds
			}

			if description != "" {
				input.Description = &description
			}

			if cycle != "" {
				input.CycleId = &cycle
			}

			if parentIssue != "" {
				input.ParentId = &parentIssue
			}

			if estimate != "" {
				n, err := strconv.Atoi(estimate)
				if err != nil {
					output.PrintErrorf("Invalid estimate: %q. Must be a number.", estimate)
				}
				input.Estimate = &n
			}

			resp, err := linear.IssueCreate(ctx, client, input)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			result := map[string]any{
				"created": resp.IssueCreate.Success,
			}
			if resp.IssueCreate.Issue != nil {
				result["id"] = resp.IssueCreate.Issue.Id
				result["identifier"] = resp.IssueCreate.Issue.Identifier
				result["url"] = resp.IssueCreate.Issue.Url
				result["title"] = resp.IssueCreate.Issue.Title
			}

			output.PrintJSON(result)
		},
	}

	cmd.Flags().StringVar(&teamFlag, "team", "", "Team ID or key")
	_ = cmd.MarkFlagRequired("team")
	cmd.Flags().StringVar(&project, "project", "", "Project ID, slug, or name")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Assignee: name, email, or user ID")
	cmd.Flags().StringVar(&priorityFlag, "priority", "", "Priority: none|urgent|high|medium|low")
	cmd.Flags().StringVar(&status, "status", "", "Status name")
	cmd.Flags().StringVar(&labels, "labels", "", "Comma-separated label names or IDs")
	cmd.Flags().StringVar(&description, "description", "", "Issue description (markdown)")
	cmd.Flags().StringVar(&cycle, "cycle", "", "Cycle ID")
	cmd.Flags().StringVar(&parentIssue, "parent", "", "Parent issue ID")
	cmd.Flags().StringVar(&estimate, "estimate", "", "Estimate points")
	parent.AddCommand(cmd)
}
