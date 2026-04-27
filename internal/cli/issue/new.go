package issue

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/priorities"
	"github.com/shhac/lin/internal/resolvers"
)

// newIssueOpts collects the flag-bound values for `lin issue new`.
type newIssueOpts struct {
	Title       string
	Team        string
	Project     string
	Assignee    string
	Priority    string
	Status      string
	Labels      string
	Description string
	Cycle       string
	ParentIssue string
	Estimate    string
}

func registerNew(parent *cobra.Command) {
	var opts newIssueOpts

	cmd := &cobra.Command{
		Use:   "new <title>",
		Short: "Create a new issue",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			opts.Title = args[0]

			client := linear.GetClient()
			input, err := buildIssueCreateInput(client, opts)
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.IssueCreate(context.Background(), client, input)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			result := map[string]any{"created": resp.IssueCreate.Success}
			if resp.IssueCreate.Issue != nil {
				result["id"] = resp.IssueCreate.Issue.Id
				result["identifier"] = resp.IssueCreate.Issue.Identifier
				result["url"] = resp.IssueCreate.Issue.Url
				result["title"] = resp.IssueCreate.Issue.Title
			}

			output.PrintJSON(result)
		},
	}

	cmd.Flags().StringVar(&opts.Team, "team", "", "Team ID or key")
	_ = cmd.MarkFlagRequired("team")
	cmd.Flags().StringVar(&opts.Project, "project", "", "Project ID, slug, or name")
	cmd.Flags().StringVar(&opts.Assignee, "assignee", "", "Assignee: name, email, or user ID")
	cmd.Flags().StringVar(&opts.Priority, "priority", "", "Priority: none|urgent|high|medium|low")
	cmd.Flags().StringVar(&opts.Status, "status", "", "Status name")
	cmd.Flags().StringVar(&opts.Labels, "labels", "", "Comma-separated label names or IDs")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Issue description (markdown)")
	cmd.Flags().StringVar(&opts.Cycle, "cycle", "", "Cycle ID")
	cmd.Flags().StringVar(&opts.ParentIssue, "parent", "", "Parent issue ID")
	cmd.Flags().StringVar(&opts.Estimate, "estimate", "", "Estimate points")
	parent.AddCommand(cmd)
}

// buildIssueCreateInput resolves human-friendly flag values to Linear IDs and
// returns the GraphQL input. Any user-input error short-circuits with a
// descriptive message; resolver errors propagate up unchanged.
func buildIssueCreateInput(client graphql.Client, opts newIssueOpts) (linear.IssueCreateInput, error) {
	priorityPtr, err := parsePriority(opts.Priority)
	if err != nil {
		return linear.IssueCreateInput{}, err
	}

	team, err := resolvers.ResolveTeam(client, opts.Team)
	if err != nil {
		return linear.IssueCreateInput{}, err
	}

	title := opts.Title
	input := linear.IssueCreateInput{
		Title:    &title,
		TeamId:   team.ID,
		Priority: priorityPtr,
	}

	if err := applyStatus(client, &input, opts.Status, team.ID); err != nil {
		return linear.IssueCreateInput{}, err
	}
	if err := applyAssignee(client, &input, opts.Assignee); err != nil {
		return linear.IssueCreateInput{}, err
	}
	if err := applyProject(client, &input, opts.Project); err != nil {
		return linear.IssueCreateInput{}, err
	}
	if err := applyLabels(client, &input, opts.Labels, team.ID); err != nil {
		return linear.IssueCreateInput{}, err
	}
	if opts.Description != "" {
		input.Description = &opts.Description
	}
	if opts.Cycle != "" {
		input.CycleId = &opts.Cycle
	}
	if opts.ParentIssue != "" {
		input.ParentId = &opts.ParentIssue
	}
	if err := applyEstimate(&input, opts.Estimate); err != nil {
		return linear.IssueCreateInput{}, err
	}
	return input, nil
}

func parsePriority(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}
	p, ok := priorities.Resolve(s)
	if !ok {
		return nil, fmt.Errorf("invalid priority: %q, valid values: %s", s, priorities.Values)
	}
	return &p, nil
}

func applyStatus(client graphql.Client, input *linear.IssueCreateInput, status, teamID string) error {
	if status == "" {
		return nil
	}
	state, err := resolvers.ResolveWorkflowState(client, status, teamID)
	if err != nil {
		return err
	}
	input.StateId = &state.ID
	return nil
}

func applyAssignee(client graphql.Client, input *linear.IssueCreateInput, assignee string) error {
	if assignee == "" {
		return nil
	}
	user, err := resolvers.ResolveUser(client, assignee)
	if err != nil {
		return err
	}
	input.AssigneeId = &user.ID
	return nil
}

func applyProject(client graphql.Client, input *linear.IssueCreateInput, project string) error {
	if project == "" {
		return nil
	}
	proj, err := resolvers.ResolveProject(client, project)
	if err != nil {
		return err
	}
	input.ProjectId = &proj.ID
	return nil
}

func applyLabels(client graphql.Client, input *linear.IssueCreateInput, labels, teamID string) error {
	if labels == "" {
		return nil
	}
	ids, err := resolvers.ResolveLabels(client, labels, teamID)
	if err != nil {
		return err
	}
	input.LabelIds = ids
	return nil
}

func applyEstimate(input *linear.IssueCreateInput, estimate string) error {
	if estimate == "" {
		return nil
	}
	n, err := strconv.Atoi(estimate)
	if err != nil {
		return fmt.Errorf("invalid estimate: %q, must be a number", estimate)
	}
	input.Estimate = &n
	return nil
}
