package issue

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	apierrors "github.com/shhac/lin/internal/errors"
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

func (f fieldUpdate) use() string {
	if f.clearable() {
		return fmt.Sprintf("%s <id> [<%s> | --clear]", f.name, f.arg)
	}
	return fmt.Sprintf("%s <id> <%s>", f.name, f.arg)
}

func registerUpdate(parent *cobra.Command) {
	update := &cobra.Command{
		Use:   "update",
		Short: "Update issue fields",
	}
	parent.AddCommand(update)

	for _, f := range fieldUpdates() {
		registerFieldUpdate(update, f)
	}

	output.HandleUnknownCommand(update, "Run 'lin issue usage' for available update subcommands")
}

func registerFieldUpdate(parent *cobra.Command, f fieldUpdate) {
	clear := false
	cmd := &cobra.Command{
		Use:   f.use(),
		Short: f.short,
		Args:  cobra.RangeArgs(1, 2),
		Run: func(_ *cobra.Command, args []string) {
			ctx := context.Background()
			if clear {
				if err := checkClearArgs(f, args); err != nil {
					output.WriteError(err)
				}
				printFieldUpdate(clearField(ctx, linear.GetClient(), f, args[0]))
				return
			}
			value, err := setValue(f, args)
			if err != nil {
				output.WriteError(err)
			}
			printFieldUpdate(setField(ctx, linear.GetClient(), f, args[0], value))
		},
	}
	if f.clearable() {
		cmd.Flags().BoolVar(&clear, "clear", false, fmt.Sprintf("Remove the issue's %s", f.arg))
	}
	parent.AddCommand(cmd)
}

func printFieldUpdate(result map[string]any, err error) {
	if err != nil {
		output.HandleGraphQLError(err)
	}
	output.PrintJSON(result)
}

func checkClearArgs(f fieldUpdate, args []string) error {
	if len(args) == 1 {
		return nil
	}
	return apierrors.Newf(apierrors.FixableByAgent, "--clear takes only the issue ID, got a %s too", f.arg).
		WithHint(fmt.Sprintf("lin issue update %s <id> --clear", f.name))
}

// setValue returns the value to set. A blank value for a clearable field is
// refused rather than sent: for a resolver-backed field Linear would match it
// against everything, and for the rest it is never what was meant.
func setValue(f fieldUpdate, args []string) (string, error) {
	if len(args) != 2 || (f.clearable() && strings.TrimSpace(args[1]) == "") {
		return "", missingValueErr(f)
	}
	return args[1], nil
}

func missingValueErr(f fieldUpdate) error {
	err := apierrors.Newf(apierrors.FixableByAgent, "missing %s", f.arg)
	if !f.clearable() {
		return err.WithHint("lin issue update " + f.use())
	}
	return err.WithHint(fmt.Sprintf("to set it: lin issue update %s <id> <%s>; to remove it: lin issue update %s <id> --clear", f.name, f.arg, f.name))
}

// setField applies the value and echoes what Linear reports afterwards, so a
// mis-resolved input is visible rather than hidden behind "updated": true.
func setField(ctx context.Context, client graphql.Client, f fieldUpdate, issueID, value string) (map[string]any, error) {
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

// clearField empties the field and confirms from Linear's response that it
// stuck; a success that left the value in place is reported as an error.
func clearField(ctx context.Context, client graphql.Client, f fieldUpdate, issueID string) (map[string]any, error) {
	resp, err := linear.IssueClear(ctx, client, issueID, f.clear)
	if err != nil {
		return nil, err
	}
	result, issue := updateResult(resp.IssueUpdate)
	if issue == nil {
		return result, nil
	}
	if _, stillSet := issue[f.key]; stillSet {
		return nil, apierrors.Newf(apierrors.FixableByHuman, "Linear reported success but the %s is still set", f.arg)
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
