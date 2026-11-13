package issue

import (
	"strings"
	"testing"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
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

func TestParseFieldUpdateArgs(t *testing.T) {
	cases := []struct {
		name      string
		field     string
		args      []string
		clear     bool
		wantValue string
		wantErr   string
		wantHint  string
	}{
		{"value", "project", []string{"ENG-1", "Roadmap"}, false, "Roadmap", "", ""},
		{"clear", "project", []string{"ENG-1"}, true, "", "", ""},
		{"clear with a value too", "project", []string{"ENG-1", "Roadmap"}, true, "", "--clear takes only the issue ID", "--clear"},
		{"missing value", "project", []string{"ENG-1"}, false, "", "missing project", "--clear"},
		{"blank value points at --clear", "project", []string{"ENG-1", ""}, false, "", "missing project", "lin issue update project <id> --clear"},
		{"whitespace value", "assignee", []string{"ENG-1", "  "}, false, "", "missing user", "lin issue update assignee <id> --clear"},
		{"blank description is a real value", "description", []string{"ENG-1", ""}, false, "", "", ""},
		{"non-clearable missing value", "title", []string{"ENG-1"}, false, "", "missing title", "lin issue update title <id> <title>"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := parseFieldUpdateArgs(fieldByName(t, tc.field), tc.args, tc.clear)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if value != tc.wantValue {
					t.Errorf("value = %q, want %q", value, tc.wantValue)
				}
				return
			}
			var apiErr *apierrors.APIError
			if !apierrors.As(err, &apiErr) {
				t.Fatalf("want APIError, got %#v", err)
			}
			if !strings.Contains(apiErr.Error(), tc.wantErr) {
				t.Errorf("error = %q, want it to contain %q", apiErr.Error(), tc.wantErr)
			}
			if !strings.Contains(apiErr.Hint, tc.wantHint) {
				t.Errorf("hint = %q, want it to contain %q", apiErr.Hint, tc.wantHint)
			}
			if !fieldByName(t, tc.field).clearable() && strings.Contains(apiErr.Hint, "--clear") {
				t.Errorf("hint = %q offers --clear on a field without it", apiErr.Hint)
			}
		})
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

func payloadWith(issue linear.IssueUpdateResultFields) linear.IssueUpdateIssueUpdateIssuePayload {
	return linear.IssueUpdateIssueUpdateIssuePayload{
		Success: true,
		Issue:   &linear.IssueUpdateIssueUpdateIssuePayloadIssue{IssueUpdateResultFields: issue},
	}
}

func TestFieldUpdateResult_EchoesResultingValue(t *testing.T) {
	issue := linear.IssueUpdateResultFields{
		Identifier: "ENG-1",
		Project:    &linear.IssueUpdateResultFieldsProject{Id: "p-uuid", Name: "Roadmap"},
	}

	got, err := fieldUpdateResult(fieldByName(t, "project"), payloadWith(issue), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	project, _ := got["project"].(map[string]any)
	if got["identifier"] != "ENG-1" || got["updated"] != true || project["name"] != "Roadmap" {
		t.Errorf("result = %#v", got)
	}
	if _, ok := got["cleared"]; ok {
		t.Errorf("a set must not report cleared: %#v", got)
	}
}

func TestFieldUpdateResult_ReportsCleared(t *testing.T) {
	got, err := fieldUpdateResult(fieldByName(t, "project"), payloadWith(linear.IssueUpdateResultFields{Identifier: "ENG-1"}), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["cleared"] != "project" || got["identifier"] != "ENG-1" {
		t.Errorf("result = %#v", got)
	}
	if _, ok := got["project"]; ok {
		t.Errorf("cleared field must not be echoed: %#v", got)
	}
}

func TestFieldUpdateResult_ClearThatDidNotStickIsAnError(t *testing.T) {
	issue := linear.IssueUpdateResultFields{
		Identifier: "ENG-1",
		Labels: linear.IssueUpdateResultFieldsLabelsIssueLabelConnection{
			Nodes: []linear.IssueUpdateResultFieldsLabelsIssueLabelConnectionNodesIssueLabel{{Id: "l-1", Name: "Bug"}},
		},
	}

	_, err := fieldUpdateResult(fieldByName(t, "labels"), payloadWith(issue), true)
	if err == nil || !strings.Contains(err.Error(), "still set") {
		t.Errorf("want a still-set error, got %v", err)
	}
}

func TestFieldUpdateResult_DescriptionIsNotEchoed(t *testing.T) {
	got, err := fieldUpdateResult(fieldByName(t, "description"), payloadWith(linear.IssueUpdateResultFields{Identifier: "ENG-1", Title: "T"}), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("want only updated + identifier, got %#v", got)
	}
}
