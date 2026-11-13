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
	// key is the MapIssueUpdateResult key echoed back; "" echoes nothing
	// (description: the caller already has it, and it can be long).
	key string
	// clear, when set, enables --clear and names the input field it empties.
	clear linear.IssueClearField
	input func(ctx context.Context, client graphql.Client, issueID, value string) (linear.IssueUpdateInput, error)
}

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
			value, err := parseFieldUpdateArgs(f, args, clear)
			if err != nil {
				output.WriteError(err)
			}
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := applyFieldUpdate(ctx, client, f, args[0], value, clear)
			if err != nil {
				output.HandleGraphQLError(err)
			}
			result, err := fieldUpdateResult(f, resp.IssueUpdate, clear)
			if err != nil {
				output.WriteError(err)
			}
			output.PrintJSON(result)
		},
	}
	if f.clearable() {
		cmd.Flags().BoolVar(&clear, "clear", false, fmt.Sprintf("Remove the issue's %s", f.arg))
	}
	parent.AddCommand(cmd)
}

// parseFieldUpdateArgs returns the value to set, or "" when clearing. A blank
// value is refused rather than sent: for a resolver-backed field Linear would
// match it against everything, and for the rest it is never what was meant.
func parseFieldUpdateArgs(f fieldUpdate, args []string, clear bool) (string, error) {
	if clear {
		if len(args) != 1 {
			return "", apierrors.Newf(apierrors.FixableByAgent, "--clear takes only the issue ID, got a %s too", f.arg).
				WithHint(fmt.Sprintf("lin issue update %s <id> --clear", f.name))
		}
		return "", nil
	}
	if len(args) != 2 {
		return "", missingValueErr(f)
	}
	if f.clearable() && strings.TrimSpace(args[1]) == "" {
		return "", missingValueErr(f)
	}
	return args[1], nil
}

func missingValueErr(f fieldUpdate) error {
	err := apierrors.Newf(apierrors.FixableByAgent, "missing %s", f.arg)
	if !f.clearable() {
		return err.WithHint(fmt.Sprintf("lin issue update %s <id> <%s>", f.name, f.arg))
	}
	return err.WithHint(fmt.Sprintf("to set it: lin issue update %s <id> <%s>; to remove it: lin issue update %s <id> --clear", f.name, f.arg, f.name))
}

func applyFieldUpdate(ctx context.Context, client graphql.Client, f fieldUpdate, issueID, value string, clear bool) (*linear.IssueUpdateResponse, error) {
	if clear {
		return linear.IssueClear(ctx, client, issueID, f.clear)
	}
	input, err := f.input(ctx, client, issueID, value)
	if err != nil {
		return nil, err
	}
	return linear.IssueUpdate(ctx, client, issueID, input)
}

// fieldUpdateResult echoes the value Linear reports after the update, so a
// mis-resolved input is visible rather than hidden behind "updated": true.
func fieldUpdateResult(f fieldUpdate, payload linear.IssueUpdateIssueUpdateIssuePayload, clear bool) (map[string]any, error) {
	result := map[string]any{"updated": payload.Success}
	if payload.Issue == nil {
		return result, nil
	}
	issue := mappers.MapIssueUpdateResult(payload.Issue.IssueUpdateResultFields)
	result["identifier"] = issue["identifier"]
	value, stillSet := issue[f.key]
	if !clear {
		if f.key != "" && stillSet {
			result[f.key] = value
		}
		return result, nil
	}
	if stillSet {
		return nil, apierrors.Newf(apierrors.FixableByHuman, "Linear reported success but the %s is still set", f.arg)
	}
	result["cleared"] = f.key
	return result, nil
}

func resolveIssueTeamID(ctx context.Context, client graphql.Client, issueID string) (string, error) {
	resp, err := linear.IssueTeam(ctx, client, issueID)
	if err != nil {
		return "", err
	}
	if resp.Issue.Team.Id == "" {
		return "", apierrors.New("could not resolve team for this issue", apierrors.FixableByAgent)
	}
	return resp.Issue.Team.Id, nil
}
