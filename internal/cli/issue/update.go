package issue

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/clearable"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
)

// fieldUpdate describes one `issue update <field>` subcommand.
type fieldUpdate struct {
	name  string // subcommand name, e.g. "due-date"
	arg   string // value placeholder in usage, e.g. "date"
	short string
	// key is the MapIssueUpdateResult key echoed back. Description leaves it
	// empty, which matches nothing: the caller already has it, and it can be long.
	key string
	// clear, when set, enables --clear and names the input field it empties.
	clear linear.IssueClearField
	input fieldInput
}

type fieldInput = func(ctx context.Context, client graphql.Client, issueID, value string) (linear.IssueUpdateInput, error)

func (f fieldUpdate) clearable() bool { return f.clear != "" }

func (f fieldUpdate) valueArg() clearable.Arg {
	return clearable.Arg{Parent: "lin issue update", Name: f.name, Placeholder: f.arg, Clearable: f.clearable()}
}

func registerUpdate(parent *cobra.Command) {
	update := &cobra.Command{
		Use:   "update",
		Short: "Update issue fields",
	}
	parent.AddCommand(update)

	for _, f := range fieldUpdates() {
		update.AddCommand(f.valueArg().Command(f.short, f.set, f.unset))
	}

	output.HandleUnknownCommand(update, "Run 'lin issue usage' for available update subcommands")
}

// set applies the value and echoes what Linear reports afterwards, so a
// mis-resolved input is visible rather than hidden behind "updated": true.
func (f fieldUpdate) set(ctx context.Context, client graphql.Client, issueID, value string) (map[string]any, error) {
	input, err := f.input(ctx, client, issueID, value)
	if err != nil {
		return nil, err
	}
	resp, err := linear.IssueUpdate(ctx, client, issueID, input)
	if err != nil {
		return nil, err
	}
	result, issue := updateResult(resp.IssueUpdate)
	if v, ok := issue[f.key]; ok {
		result[f.key] = v
	}
	return result, nil
}

// unset empties the field and confirms from Linear's response that it
// stuck; a success that left the value in place is reported as an error.
func (f fieldUpdate) unset(ctx context.Context, client graphql.Client, issueID string) (map[string]any, error) {
	resp, err := linear.IssueClear(ctx, client, issueID, f.clear)
	if err != nil {
		return nil, err
	}
	result, issue := updateResult(resp.IssueUpdate)
	if issue == nil {
		return result, nil
	}
	if _, stillSet := issue[f.key]; stillSet {
		return nil, f.valueArg().StillSetErr()
	}
	result["cleared"] = f.key
	return result, nil
}

// updateResult is the shared base of every update's output, plus the issue
// as Linear reports it (nil when the payload carries none).
func updateResult(payload linear.IssueUpdateIssueUpdateIssuePayload) (map[string]any, map[string]any) {
	result := map[string]any{"updated": payload.Success}
	if payload.Issue == nil {
		return result, nil
	}
	issue := mappers.MapIssueUpdateResult(payload.Issue.IssueUpdateResultFields)
	result["identifier"] = issue["identifier"]
	return result, issue
}
