package clearable

import (
	"context"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"

	apierrors "github.com/shhac/lin/internal/errors"
)

var (
	lead  = Arg{Parent: "lin project update", Name: "lead", Placeholder: "user", Clearable: true}
	title = Arg{Parent: "lin project update", Name: "title", Placeholder: "title"}
)

func assertAPIError(t *testing.T, err error, wantMsg, wantHint string) *apierrors.APIError {
	t.Helper()
	var apiErr *apierrors.APIError
	if !apierrors.As(err, &apiErr) {
		t.Fatalf("want APIError, got %#v", err)
	}
	if !strings.Contains(apiErr.Error(), wantMsg) {
		t.Errorf("error = %q, want it to contain %q", apiErr.Error(), wantMsg)
	}
	if !strings.Contains(apiErr.Hint, wantHint) {
		t.Errorf("hint = %q, want it to contain %q", apiErr.Hint, wantHint)
	}
	return apiErr
}

func TestUse(t *testing.T) {
	if got := lead.use(); got != "lead <id> [<user> | --clear]" {
		t.Errorf("clearable Use = %q", got)
	}
	if got := title.use(); got != "title <id> <title>" {
		t.Errorf("plain Use = %q", got)
	}
}

func TestValue(t *testing.T) {
	cases := []struct {
		name     string
		arg      Arg
		args     []string
		want     string
		wantErr  string
		wantHint string
	}{
		{"value", lead, []string{"P-1", "Alice"}, "Alice", "", ""},
		{"missing value", lead, []string{"P-1"}, "", "missing user", "lin project update lead <id> --clear"},
		{"blank value points at --clear", lead, []string{"P-1", ""}, "", "missing user", "to remove it: lin project update lead <id> --clear"},
		{"whitespace value", lead, []string{"P-1", " \t"}, "", "missing user", "--clear"},
		{"blank is a real value when not clearable", title, []string{"P-1", ""}, "", "", ""},
		{"missing value without --clear", title, []string{"P-1"}, "", "missing title", "lin project update title <id> <title>"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.arg.Value(tc.args)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if value != tc.want {
					t.Errorf("value = %q, want %q", value, tc.want)
				}
				return
			}
			apiErr := assertAPIError(t, err, tc.wantErr, tc.wantHint)
			if !tc.arg.Clearable && strings.Contains(apiErr.Hint, "--clear") {
				t.Errorf("hint = %q offers --clear on an argument without it", apiErr.Hint)
			}
		})
	}
}

func TestCheckClearArgs(t *testing.T) {
	if err := lead.checkClearArgs([]string{"P-1"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	err := lead.checkClearArgs([]string{"P-1", "Alice"})
	assertAPIError(t, err, "--clear takes only the ID; drop the user value", "lin project update lead <id> --clear")
}

func TestStillSetErr(t *testing.T) {
	err := Arg{Name: "due-date", Clearable: true}.StillSetErr()
	apiErr := assertAPIError(t, err, "did not remove the due date", "")
	if apiErr.FixableBy != apierrors.FixableByHuman {
		t.Errorf("fixable_by = %q, want human", apiErr.FixableBy)
	}
}

func TestCommand_OnlyClearableArgsGetTheFlag(t *testing.T) {
	if lead.Command("s", nil, nil).Flags().Lookup("clear") == nil {
		t.Error("clearable command is missing --clear")
	}
	if title.Command("s", nil, nil).Flags().Lookup("clear") != nil {
		t.Error("non-clearable command must not offer --clear")
	}
}

func TestRun_DispatchesAfterValidatingWithoutAClient(t *testing.T) {
	var calls []string
	set := func(_ context.Context, _ graphql.Client, id, value string) (map[string]any, error) {
		calls = append(calls, "set "+id+" "+value)
		return map[string]any{"set": value}, nil
	}
	clear := func(_ context.Context, _ graphql.Client, id string) (map[string]any, error) {
		calls = append(calls, "clear "+id)
		return map[string]any{"cleared": id}, nil
	}
	cases := []struct {
		name      string
		args      []string
		clearing  bool
		wantCalls string
		wantErr   string
	}{
		{"set", []string{"P-1", "Alice"}, false, "set P-1 Alice", ""},
		{"clear", []string{"P-1"}, true, "clear P-1", ""},
		{"clear with a value is refused", []string{"P-1", "Alice"}, true, "", "--clear takes only the ID"},
		{"missing value is refused", []string{"P-1"}, false, "", "missing user"},
		{"blank value is refused", []string{"P-1", " "}, false, "", "missing user"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls = nil
			clients := 0
			newClient := func() graphql.Client { clients++; return nil }

			_, err := lead.run(context.Background(), tc.args, tc.clearing, newClient, set, clear)

			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("err = %v, want %q", err, tc.wantErr)
				}
				if clients != 0 || len(calls) != 0 {
					t.Errorf("a usage error must not create a client or run (clients=%d, calls=%v)", clients, calls)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.Join(calls, ";") != tc.wantCalls || clients != 1 {
				t.Errorf("calls = %v, clients = %d, want %q once", calls, clients, tc.wantCalls)
			}
		})
	}
}
