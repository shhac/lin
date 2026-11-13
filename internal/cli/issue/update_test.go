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

func TestFieldUpdateArgs(t *testing.T) {
	project := fieldByName(t, "project").valueArg()
	if project.Parent != "lin issue update" || project.Name != "project" || !project.Clearable {
		t.Errorf("project arg = %+v", project)
	}
	if _, err := fieldByName(t, "description").valueArg().Value([]string{"ENG-1", ""}); err != nil {
		t.Errorf("a blank description is a real value, got %v", err)
	}
	if _, err := project.Value([]string{"ENG-1", ""}); err == nil {
		t.Error("a blank project must be refused")
	}
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

	got, err := fieldByName(t, "title").set(context.Background(), mock.Client(), "ENG-1", "New title")
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

	got, err := fieldByName(t, "description").set(context.Background(), mock.Client(), "ENG-1", "")
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

			if _, err := f.set(context.Background(), mock.Client(), "ENG-1", value); err != nil {
				t.Fatalf("set: %v", err)
			}
			setInput, _ := mock.Variables("IssueUpdate")["input"].(map[string]any)
			if _, ok := setInput[string(f.clear)]; !ok || len(setInput) != 1 {
				t.Errorf("set sent %v, want only %q", setInput, f.clear)
			}

			if _, err := f.unset(context.Background(), mock.Client(), "ENG-1"); err != nil {
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

	got, err := fieldByName(t, "project").unset(context.Background(), mock.Client(), "ENG-1")
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

	_, err := fieldByName(t, "labels").unset(context.Background(), mock.Client(), "ENG-1")
	if err == nil || !strings.Contains(err.Error(), "did not remove the labels") {
		t.Errorf("want a not-removed error, got %v", err)
	}
}

func TestClearField_WithoutIssueInPayloadDoesNotClaimCleared(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	mock.Handle("IssueUpdate", issueUpdateResponse(nil))

	got, err := fieldByName(t, "project").unset(context.Background(), mock.Client(), "ENG-1")
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
	var apiErr *apierrors.APIError
	if !apierrors.As(err, &apiErr) || !strings.Contains(apiErr.Error(), "invalid estimate") {
		t.Errorf("want an invalid-estimate APIError, got %#v", err)
	}

	_, err = estimateInput(ctx, mock.Client(), "ENG-1", "4")
	if !apierrors.As(err, &apiErr) {
		t.Errorf("want an off-scale APIError, got %#v", err)
	}
}
