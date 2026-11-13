package issue

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/priorities"
	"github.com/shhac/lin/internal/resolvers"
)

func fieldUpdates() []fieldUpdate {
	return []fieldUpdate{
		{name: "status", arg: "status", short: "Update issue status", key: "status", input: statusInput},
		{name: "assignee", arg: "user", short: "Update or remove issue assignee", key: "assignee", clear: linear.ClearAssignee, input: assigneeInput},
		{name: "priority", arg: "priority", short: "Update issue priority", key: "priorityLabel", input: priorityInput},
		{name: "project", arg: "project", short: "Move issue to a project, or remove it from its project", key: "project", clear: linear.ClearProject, input: projectInput},
		{name: "labels", arg: "labels", short: "Replace issue labels, or remove them all", key: "labels", clear: linear.ClearLabels, input: labelsInput},
		{name: "estimate", arg: "estimate", short: "Update issue estimate (validated against team scale), or remove it", key: "estimate", clear: linear.ClearEstimate, input: estimateInput},
		{name: "title", arg: "title", short: "Update issue title", key: "title", input: literalInput(func(v string) linear.IssueUpdateInput { return linear.IssueUpdateInput{Title: &v} })},
		{name: "description", arg: "description", short: "Update issue description", input: literalInput(func(v string) linear.IssueUpdateInput { return linear.IssueUpdateInput{Description: &v} })},
		{name: "due-date", arg: "date", short: "Update or remove issue due date", key: "dueDate", clear: linear.ClearDueDate, input: literalInput(func(v string) linear.IssueUpdateInput { return linear.IssueUpdateInput{DueDate: &v} })},
		{name: "cycle", arg: "cycle-id", short: "Move issue to a cycle, or remove it from its cycle", key: "cycle", clear: linear.ClearCycle, input: literalInput(func(v string) linear.IssueUpdateInput { return linear.IssueUpdateInput{CycleId: &v} })},
		{name: "parent", arg: "parent-id", short: "Set parent issue (make sub-issue), or detach from parent", key: "parent", clear: linear.ClearParent, input: literalInput(func(v string) linear.IssueUpdateInput { return linear.IssueUpdateInput{ParentId: &v} })},
	}
}

type fieldInput = func(ctx context.Context, client graphql.Client, issueID, value string) (linear.IssueUpdateInput, error)

func literalInput(build func(string) linear.IssueUpdateInput) fieldInput {
	return func(_ context.Context, _ graphql.Client, _, value string) (linear.IssueUpdateInput, error) {
		return build(value), nil
	}
}

func statusInput(ctx context.Context, client graphql.Client, issueID, value string) (linear.IssueUpdateInput, error) {
	teamID, err := resolveIssueTeamID(ctx, client, issueID)
	if err != nil {
		return linear.IssueUpdateInput{}, err
	}
	state, err := resolvers.ResolveWorkflowState(client, value, teamID)
	if err != nil {
		return linear.IssueUpdateInput{}, err
	}
	return linear.IssueUpdateInput{StateId: &state.ID}, nil
}

func assigneeInput(_ context.Context, client graphql.Client, _, value string) (linear.IssueUpdateInput, error) {
	user, err := resolvers.ResolveUser(client, value)
	if err != nil {
		return linear.IssueUpdateInput{}, err
	}
	return linear.IssueUpdateInput{AssigneeId: &user.ID}, nil
}

func priorityInput(_ context.Context, _ graphql.Client, _, value string) (linear.IssueUpdateInput, error) {
	p, ok := priorities.Resolve(value)
	if !ok {
		return linear.IssueUpdateInput{}, apierrors.Newf(apierrors.FixableByAgent, "invalid priority: %q, valid values: %s", value, priorities.Values)
	}
	return linear.IssueUpdateInput{Priority: &p}, nil
}

func projectInput(_ context.Context, client graphql.Client, _, value string) (linear.IssueUpdateInput, error) {
	proj, err := resolvers.ResolveProject(client, value)
	if err != nil {
		return linear.IssueUpdateInput{}, err
	}
	return linear.IssueUpdateInput{ProjectId: &proj.ID}, nil
}

func labelsInput(ctx context.Context, client graphql.Client, issueID, value string) (linear.IssueUpdateInput, error) {
	teamID, err := resolveIssueTeamID(ctx, client, issueID)
	if err != nil {
		return linear.IssueUpdateInput{}, err
	}
	ids, err := resolvers.ResolveIssueLabels(client, value, teamID)
	if err != nil {
		return linear.IssueUpdateInput{}, err
	}
	return linear.IssueUpdateInput{LabelIds: ids}, nil
}
