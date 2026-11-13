package issue

import (
	"context"
	"reflect"
	"strings"
	"testing"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/testutil"
)

func fieldByName(t *testing.T, name string) fieldUpdate {
	t.Helper()
	for _, f := range fieldUpdates() {
		if f.name == name {
			return f
		}
	}
	t.Fatalf("no field update named %q", name)
	return fieldUpdate{}
}

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

func TestSetValue(t *testing.T) {
	cases := []struct {
		name     string
		field    string
		args     []string
		want     string
		wantErr  string
		wantHint string
	}{
		{"value", "project", []string{"ENG-1", "Roadmap"}, "Roadmap", "", ""},
		{"missing value", "project", []string{"ENG-1"}, "", "missing project", "--clear"},
		{"blank value points at --clear", "project", []string{"ENG-1", ""}, "", "missing project", "lin issue update project <id> --clear"},
		{"whitespace value", "assignee", []string{"ENG-1", "  "}, "", "missing user", "lin issue update assignee <id> --clear"},
		{"blank description is a real value", "description", []string{"ENG-1", ""}, "", "", ""},
		{"non-clearable missing value", "title", []string{"ENG-1"}, "", "missing title", "lin issue update title <id> <title>"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := fieldByName(t, tc.field)
			value, err := setValue(f, tc.args)
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
			if !f.clearable() && strings.Contains(apiErr.Hint, "--clear") {
				t.Errorf("hint = %q offers --clear on a field without it", apiErr.Hint)
			}
		})
	}
}

func TestCheckClearArgs(t *testing.T) {
	f := fieldByName(t, "project")
	if err := checkClearArgs(f, []string{"ENG-1"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	err := checkClearArgs(f, []string{"ENG-1", "Roadmap"})
	assertAPIError(t, err, "--clear takes only the issue ID", "lin issue update project <id> --clear")
}

func TestOnlyNullableFieldsAreClearable(t *testing.T) {
	clearable := map[string]bool{}
	for _, f := range fieldUpdates() {
		if f.clearable() {
			clearable[f.name] = true
		}
	}
	for _, name := range []string{"project", "assignee", "cycle", "parent", "due-date", "estimate", "labels"} {
		if !clearable[name] {
			t.Errorf("%s should support --clear", name)
		}
	}
	// Linear rejects null for these with "Argument Validation Error".
	for _, name := range []string{"status", "title", "priority", "description"} {
		if clearable[name] {
			t.Errorf("%s must not offer --clear", name)
		}
	}
}

// fullyPopulatedIssue sets every field MapIssueUpdateResult can emit.
func fullyPopulatedIssue() linear.IssueUpdateResultFields {
	estimate, dueDate, cycleName := 3.0, "2026-12-31", "Cycle 9"
	return linear.IssueUpdateResultFields{
		Identifier:    "ENG-1",
		Title:         "T",
		PriorityLabel: "High",
		Estimate:      &estimate,
		DueDate:       &dueDate,
		State:         linear.IssueUpdateResultFieldsStateWorkflowState{Id: "s-1", Name: "Todo", Type: "unstarted"},
		Assignee:      &linear.IssueUpdateResultFieldsAssigneeUser{Id: "u-1", Name: "Test User"},
		Project:       &linear.IssueUpdateResultFieldsProject{Id: "p-1", Name: "Roadmap"},
		Cycle:         &linear.IssueUpdateResultFieldsCycle{Id: "c-1", Number: 9, Name: &cycleName},
		Parent:        &linear.IssueUpdateResultFieldsParentIssue{Id: "i-0", Identifier: "ENG-0"},
		Labels: linear.IssueUpdateResultFieldsLabelsIssueLabelConnection{
			Nodes: []linear.IssueUpdateResultFieldsLabelsIssueLabelConnectionNodesIssueLabel{{Id: "l-1", Name: "Bug"}},
		},
	}
}

// The --clear "still set" check and the echo both look the spec's key up in
// the mapper's output, so a key the mapper never emits would silently pass
// every clear and echo nothing.
func TestFieldKeysMatchMapperOutput(t *testing.T) {
	populated := mappers.MapIssueUpdateResult(fullyPopulatedIssue())
	empty := mappers.MapIssueUpdateResult(linear.IssueUpdateResultFields{Identifier: "ENG-1"})
	for _, f := range fieldUpdates() {
		if f.key == "" {
			if f.clearable() {
				t.Errorf("%s is clearable but has no key to verify the clear against", f.name)
			}
			continue
		}
		if _, ok := populated[f.key]; !ok {
			t.Errorf("%s: key %q is not emitted by MapIssueUpdateResult", f.name, f.key)
		}
		if _, ok := empty[f.key]; ok && f.clearable() {
			t.Errorf("%s: key %q is emitted even when unset, so a clear could never be confirmed", f.name, f.key)
		}
	}
}

func issueUpdateResponse(issue *linear.IssueUpdateResultFields) map[string]any {
	payload := map[string]any{"success": true}
	if issue != nil {
		payload["issue"] = issue
	}
	return map[string]any{"issueUpdate": payload}
}

func TestSetField_EchoesResultingValue(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	issue := linear.IssueUpdateResultFields{Identifier: "ENG-1", Title: "New title"}
	mock.Handle("IssueUpdate", issueUpdateResponse(&issue))

	got, err := setField(context.Background(), mock.Client(), fieldByName(t, "title"), "ENG-1", "New title")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]any{"updated": true, "identifier": "ENG-1", "title": "New title"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestSetField_DescriptionIsNotEchoed(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	issue := linear.IssueUpdateResultFields{Identifier: "ENG-1", Title: "T"}
	mock.Handle("IssueUpdate", issueUpdateResponse(&issue))

	got, err := setField(context.Background(), mock.Client(), fieldByName(t, "description"), "ENG-1", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := map[string]any{"updated": true, "identifier": "ENG-1"}; !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

// Each clearable field must send the same input key whether it is being set
// or cleared — a spec row pairing one field's input with another's clear
// constant would empty the wrong field.
func TestClearableFields_SetAndClearTargetTheSameInput(t *testing.T) {
	literal := map[string]string{"due-date": "2026-12-31", "cycle": "c-uuid", "parent": "i-uuid"}
	for _, f := range fieldUpdates() {
		value, ok := literal[f.name]
		if !ok {
			continue
		}
		t.Run(f.name, func(t *testing.T) {
			mock := testutil.NewMockLinear()
			defer mock.Close()
			mock.Handle("IssueUpdate", issueUpdateResponse(&linear.IssueUpdateResultFields{Identifier: "ENG-1"}))

			if _, err := setField(context.Background(), mock.Client(), f, "ENG-1", value); err != nil {
				t.Fatalf("set: %v", err)
			}
			setInput, _ := mock.Variables("IssueUpdate")["input"].(map[string]any)
			if _, ok := setInput[string(f.clear)]; !ok || len(setInput) != 1 {
				t.Errorf("set sent %v, want only %q", setInput, f.clear)
			}

			if _, err := clearField(context.Background(), mock.Client(), f, "ENG-1"); err != nil {
				t.Fatalf("clear: %v", err)
			}
			clearInput, _ := mock.Variables("IssueUpdate")["input"].(map[string]any)
			if _, ok := clearInput[string(f.clear)]; !ok || len(clearInput) != 1 {
				t.Errorf("clear sent %v, want only %q", clearInput, f.clear)
			}
		})
	}
}

func TestClearField_ReportsCleared(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	mock.Handle("IssueUpdate", issueUpdateResponse(&linear.IssueUpdateResultFields{Identifier: "ENG-1"}))

	got, err := clearField(context.Background(), mock.Client(), fieldByName(t, "project"), "ENG-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := map[string]any{"updated": true, "identifier": "ENG-1", "cleared": "project"}; !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestClearField_ThatDidNotStickIsAnError(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	issue := fullyPopulatedIssue()
	mock.Handle("IssueUpdate", issueUpdateResponse(&issue))

	_, err := clearField(context.Background(), mock.Client(), fieldByName(t, "labels"), "ENG-1")
	if err == nil || !strings.Contains(err.Error(), "still set") {
		t.Errorf("want a still-set error, got %v", err)
	}
}

func TestClearField_WithoutIssueInPayloadDoesNotClaimCleared(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	mock.Handle("IssueUpdate", issueUpdateResponse(nil))

	got, err := clearField(context.Background(), mock.Client(), fieldByName(t, "project"), "ENG-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := map[string]any{"updated": true}; !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestEstimateInput(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	mock.Handle("IssueTeam", map[string]any{"issue": map[string]any{"team": map[string]any{"id": "t-1"}}})
	mock.Handle("TeamGet", map[string]any{"team": map[string]any{
		"id": "t-1", "name": "Eng", "key": "ENG",
		"issueEstimationType": "fibonacci", "issueEstimationAllowZero": false, "issueEstimationExtended": false,
	}})
	ctx := context.Background()

	input, err := estimateInput(ctx, mock.Client(), "ENG-1", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Estimate == nil || *input.Estimate != 5 {
		t.Errorf("estimate = %v, want 5", input.Estimate)
	}

	_, err = estimateInput(ctx, mock.Client(), "ENG-1", "five")
	assertAPIError(t, err, "invalid estimate", "")

	_, err = estimateInput(ctx, mock.Client(), "ENG-1", "4")
	assertAPIError(t, err, "", "")
}
