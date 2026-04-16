package mappers

import "testing"

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
