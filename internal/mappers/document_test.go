package mappers

import "testing"

func TestMapDocSummary_Full(t *testing.T) {
	input := DocSummaryInput{
		ID:          "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		SlugId:      "doc-abc",
		Title:       "Architecture Overview",
		URL:         "https://linear.app/test/document/doc-abc",
		UpdatedAt:   "2025-03-15T10:30:00.000Z",
		CreatorID:   "11111111-2222-3333-4444-555555555555",
		CreatorName: "Alan Turing",
		ProjectID:   "22222222-3333-4444-5555-666666666666",
		ProjectName: "Platform Migration",
	}
	got := MapDocSummary(input)

	if got["id"] != "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" {
		t.Errorf("id = %v", got["id"])
	}
	if got["slugId"] != "doc-abc" {
		t.Errorf("slugId = %v", got["slugId"])
	}
	if got["title"] != "Architecture Overview" {
		t.Errorf("title = %v", got["title"])
	}
	if got["updatedAt"] != "2025-03-15T10:30:00.000Z" {
		t.Errorf("updatedAt = %v", got["updatedAt"])
	}

	creator := got["creator"].(map[string]any)
	if creator["id"] != "11111111-2222-3333-4444-555555555555" {
		t.Errorf("creator.id = %v", creator["id"])
	}
	if creator["name"] != "Alan Turing" {
		t.Errorf("creator.name = %v", creator["name"])
	}

	project := got["project"].(map[string]any)
	if project["id"] != "22222222-3333-4444-5555-666666666666" {
		t.Errorf("project.id = %v", project["id"])
	}
	if project["name"] != "Platform Migration" {
		t.Errorf("project.name = %v", project["name"])
	}
}

func TestMapDocSummary_NoOptionalFields(t *testing.T) {
	input := DocSummaryInput{
		ID:        "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		SlugId:    "doc-xyz",
		Title:     "Standalone Doc",
		URL:       "https://linear.app/test/document/doc-xyz",
		UpdatedAt: "2025-04-01T08:00:00.000Z",
	}
	got := MapDocSummary(input)

	if _, ok := got["creator"]; ok {
		t.Error("creator should be absent when empty")
	}
	if _, ok := got["project"]; ok {
		t.Error("project should be absent when empty")
	}
}
