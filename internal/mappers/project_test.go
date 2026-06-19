package mappers

import (
	"testing"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/ptr"
)

func TestMapProjectSummary_Full(t *testing.T) {
	input := ProjectSummaryInput{
		ID:         "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		SlugId:     "proj-abc",
		URL:        "https://linear.app/test/project/proj-abc",
		Name:       "Platform Migration",
		State:      "started",
		Progress:   0.65,
		LeadName:   "Grace Hopper",
		StartDate:  "2025-01-15",
		TargetDate: "2025-06-30",
	}
	got := MapProjectSummary(input)

	if got["id"] != "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" {
		t.Errorf("id = %v", got["id"])
	}
	if got["slugId"] != "proj-abc" {
		t.Errorf("slugId = %v", got["slugId"])
	}
	if got["name"] != "Platform Migration" {
		t.Errorf("name = %v", got["name"])
	}
	if got["status"] != "started" {
		t.Errorf("status = %v", got["status"])
	}
	if got["progress"] != 0.65 {
		t.Errorf("progress = %v", got["progress"])
	}
	if got["lead"] != "Grace Hopper" {
		t.Errorf("lead = %v", got["lead"])
	}
	if got["startDate"] != "2025-01-15" {
		t.Errorf("startDate = %v", got["startDate"])
	}
	if got["targetDate"] != "2025-06-30" {
		t.Errorf("targetDate = %v", got["targetDate"])
	}
}

func TestFromProjectUpdateSummary(t *testing.T) {
	f := linear.ProjectUpdateSummaryFields{
		Id:        "uuuuuuuu-1111-2222-3333-444444444444",
		Url:       "https://linear.app/test/projectUpdate/abc",
		Health:    linear.ProjectUpdateHealthTypeOntrack,
		Body:      "Shipped the thing",
		CreatedAt: "2026-06-19T10:00:00.000Z",
		EditedAt:  ptr.To("2026-06-19T11:00:00.000Z"),
		User: linear.ProjectUpdateSummaryFieldsUser{
			Id:   "user-1",
			Name: "Grace Hopper",
		},
	}
	got := FromProjectUpdateSummary(f)

	if got["id"] != f.Id {
		t.Errorf("id = %v", got["id"])
	}
	if got["health"] != "onTrack" {
		t.Errorf("health = %v, want onTrack (string)", got["health"])
	}
	if got["body"] != "Shipped the thing" {
		t.Errorf("body = %v", got["body"])
	}
	if got["editedAt"] != "2026-06-19T11:00:00.000Z" {
		t.Errorf("editedAt = %v", got["editedAt"])
	}
	user, ok := got["user"].(map[string]any)
	if !ok || user["name"] != "Grace Hopper" {
		t.Errorf("user = %v", got["user"])
	}
}

func TestFromProjectUpdateSummary_NoEditedAt(t *testing.T) {
	f := linear.ProjectUpdateSummaryFields{
		Id:        "uuuuuuuu-1111-2222-3333-444444444444",
		Health:    linear.ProjectUpdateHealthTypeAtrisk,
		CreatedAt: "2026-06-19T10:00:00.000Z",
		User:      linear.ProjectUpdateSummaryFieldsUser{Id: "user-1", Name: "Ada"},
	}
	got := FromProjectUpdateSummary(f)
	if _, ok := got["editedAt"]; ok {
		t.Error("editedAt should be absent when nil")
	}
}

func TestMapProjectSummary_NoOptionalFields(t *testing.T) {
	input := ProjectSummaryInput{
		ID:       "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		SlugId:   "proj-xyz",
		URL:      "https://linear.app/test/project/proj-xyz",
		Name:     "Empty Project",
		State:    "backlog",
		Progress: 0,
	}
	got := MapProjectSummary(input)

	if _, ok := got["lead"]; ok {
		t.Error("lead should be absent when empty")
	}
	if _, ok := got["startDate"]; ok {
		t.Error("startDate should be absent when empty")
	}
	if _, ok := got["targetDate"]; ok {
		t.Error("targetDate should be absent when empty")
	}
}
