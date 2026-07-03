package initiative

import (
	"strings"
	"testing"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/testutil"
)

func TestBuildInitiativeCreateInput_RequiredOnly(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	input, err := buildInitiativeCreateInput(mock.Client(), newInitiativeOpts{
		Name: "Roadmap",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Name != "Roadmap" {
		t.Errorf("name = %q", input.Name)
	}
	if input.Description != nil || input.Status != nil || input.OwnerId != nil {
		t.Error("optional fields should be unset")
	}
}

func TestBuildInitiativeCreateInput_Optionals(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	input, err := buildInitiativeCreateInput(mock.Client(), newInitiativeOpts{
		Name:        "Roadmap",
		Description: "details",
		Content:     "body",
		Color:       "#fff",
		Icon:        "Rocket",
		TargetDate:  "2026-03-01",
		Status:      "active",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Description == nil || *input.Description != "details" {
		t.Errorf("description = %v", input.Description)
	}
	if input.Content == nil || *input.Content != "body" {
		t.Errorf("content = %v", input.Content)
	}
	if input.Color == nil || *input.Color != "#fff" {
		t.Errorf("color = %v", input.Color)
	}
	if input.Icon == nil || *input.Icon != "Rocket" {
		t.Errorf("icon = %v", input.Icon)
	}
	if input.TargetDate == nil || *input.TargetDate != "2026-03-01" {
		t.Errorf("targetDate = %v", input.TargetDate)
	}
	if input.Status == nil || *input.Status != linear.InitiativeStatus("Active") {
		t.Errorf("status should be capitalized to Active, got %v", input.Status)
	}
}

func TestBuildInitiativeCreateInput_InvalidStatus(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	_, err := buildInitiativeCreateInput(mock.Client(), newInitiativeOpts{
		Name:   "Roadmap",
		Status: "paused",
	})
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
	if !strings.Contains(err.Error(), "unknown initiative status") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuildInitiativeCreateInput_OwnerResolution(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("UserList", map[string]any{
		"users": map[string]any{
			"nodes": []map[string]any{
				{"id": "user-uuid", "name": "Ada Lovelace", "email": "ada@example.com", "displayName": "ada"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	input, err := buildInitiativeCreateInput(mock.Client(), newInitiativeOpts{
		Name:  "Roadmap",
		Owner: "ada@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.OwnerId == nil || *input.OwnerId != "user-uuid" {
		t.Errorf("ownerId = %v", input.OwnerId)
	}
}
