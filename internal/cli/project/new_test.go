package project

import (
	"strings"
	"testing"

	"github.com/shhac/lin/internal/testutil"
)

func handleTeamGet(mock *testutil.MockLinear) {
	mock.Handle("TeamGet", map[string]any{
		"team": map[string]any{"id": "team-uuid", "name": "Engineering", "key": "ENG"},
	})
}

func TestBuildProjectCreateInput_RequiredOnly(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	handleTeamGet(mock)

	input, err := buildProjectCreateInput(mock.Client(), newProjectOpts{
		Name: "Launch",
		Team: "ENG",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Name != "Launch" {
		t.Errorf("name = %q", input.Name)
	}
	if len(input.TeamIds) != 1 || input.TeamIds[0] != "team-uuid" {
		t.Errorf("teamIds = %v", input.TeamIds)
	}
	if input.Description != nil || input.LeadId != nil || input.State != nil {
		t.Error("optional fields should be unset")
	}
}

func TestBuildProjectCreateInput_Optionals(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	handleTeamGet(mock)

	input, err := buildProjectCreateInput(mock.Client(), newProjectOpts{
		Name:        "Launch",
		Team:        "ENG",
		Description: "details",
		StartDate:   "2026-01-01",
		TargetDate:  "2026-02-01",
		Status:      "STARTED",
		Content:     "body",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Description == nil || *input.Description != "details" {
		t.Errorf("description = %v", input.Description)
	}
	if input.StartDate == nil || *input.StartDate != "2026-01-01" {
		t.Errorf("startDate = %v", input.StartDate)
	}
	if input.TargetDate == nil || *input.TargetDate != "2026-02-01" {
		t.Errorf("targetDate = %v", input.TargetDate)
	}
	if input.State == nil || *input.State != "started" {
		t.Errorf("state should be normalized to lowercase, got %v", input.State)
	}
	if input.Content == nil || *input.Content != "body" {
		t.Errorf("content = %v", input.Content)
	}
}

func TestBuildProjectCreateInput_InvalidStatus(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	handleTeamGet(mock)

	_, err := buildProjectCreateInput(mock.Client(), newProjectOpts{
		Name:   "Launch",
		Team:   "ENG",
		Status: "shipping",
	})
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
	if !strings.Contains(err.Error(), "unknown project status") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestResolveTeamIDs_SkipsBlanks(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	handleTeamGet(mock)

	ids, err := resolveTeamIDs(mock.Client(), "ENG, ,ENG")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 ids (blank entry skipped), got %v", ids)
	}
}

func TestResolveTeamIDs_Empty(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	ids, err := resolveTeamIDs(mock.Client(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected no ids for empty input, got %v", ids)
	}
}
