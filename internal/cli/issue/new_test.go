package issue

import (
	"strings"
	"testing"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/testutil"
)

func TestParsePriority_ValidValues(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"none", 0},
		{"urgent", 1},
		{"high", 2},
		{"medium", 3},
		{"low", 4},
		{"HIGH", 2},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, err := parsePriority(tc.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("expected non-nil pointer")
			}
			if *got != tc.want {
				t.Errorf("priority = %d, want %d", *got, tc.want)
			}
		})
	}
}

func TestParsePriority_EmptyReturnsNil(t *testing.T) {
	got, err := parsePriority("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil pointer for empty input, got %d", *got)
	}
}

func TestParsePriority_Invalid(t *testing.T) {
	_, err := parsePriority("critical")
	if err == nil {
		t.Fatal("expected error for invalid priority")
	}
	if !strings.Contains(err.Error(), "invalid priority") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestApplyEstimate_ValidNumber(t *testing.T) {
	input := linear.IssueCreateInput{}
	if err := applyEstimate(&input, "5"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Estimate == nil || *input.Estimate != 5 {
		t.Errorf("estimate = %v", input.Estimate)
	}
}

func TestApplyEstimate_EmptySkips(t *testing.T) {
	input := linear.IssueCreateInput{}
	if err := applyEstimate(&input, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Estimate != nil {
		t.Error("expected estimate to remain nil when input is empty")
	}
}

func TestApplyEstimate_NonNumeric(t *testing.T) {
	input := linear.IssueCreateInput{}
	if err := applyEstimate(&input, "five"); err == nil {
		t.Fatal("expected error for non-numeric estimate")
	}
}

func TestBuildIssueCreateInput_RequiredOnly(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("TeamGet", map[string]any{
		"team": map[string]any{
			"id":   "team-uuid",
			"name": "Engineering",
			"key":  "ENG",
		},
	})

	input, err := buildIssueCreateInput(mock.Client(), newIssueOpts{
		Title: "My new issue",
		Team:  "ENG",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Title == nil || *input.Title != "My new issue" {
		t.Errorf("title = %v", input.Title)
	}
	if input.TeamId != "team-uuid" {
		t.Errorf("teamId = %q", input.TeamId)
	}
	if input.Priority != nil {
		t.Errorf("priority should be nil when unset, got %v", *input.Priority)
	}
	if input.AssigneeId != nil || input.ProjectId != nil || input.StateId != nil {
		t.Error("optional fields should be unset")
	}
}

func TestBuildIssueCreateInput_PriorityAndDescription(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("TeamGet", map[string]any{
		"team": map[string]any{"id": "team-uuid", "name": "Engineering", "key": "ENG"},
	})

	input, err := buildIssueCreateInput(mock.Client(), newIssueOpts{
		Title:       "Title",
		Team:        "ENG",
		Priority:    "urgent",
		Description: "details",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Priority == nil || *input.Priority != 1 {
		t.Errorf("priority = %v", input.Priority)
	}
	if input.Description == nil || *input.Description != "details" {
		t.Errorf("description = %v", input.Description)
	}
}

func TestBuildIssueCreateInput_InvalidPriority(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	_, err := buildIssueCreateInput(mock.Client(), newIssueOpts{
		Title:    "T",
		Team:     "ENG",
		Priority: "critical",
	})
	if err == nil {
		t.Fatal("expected error for invalid priority")
	}
	if !strings.Contains(err.Error(), "invalid priority") {
		t.Errorf("unexpected error: %v", err)
	}
}

