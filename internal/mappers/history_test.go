package mappers

import (
	"testing"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/ptr"
)

func TestMapHistoryEntry_StatusChange(t *testing.T) {
	h := HistoryNode{
		Id:        "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		CreatedAt: "2025-04-01T10:00:00.000Z",
		Actor: &linear.IssueHistoryIssueHistoryIssueHistoryConnectionNodesIssueHistoryActorUser{
			Id:   "11111111-2222-3333-4444-555555555555",
			Name: "Ada Lovelace",
		},
		FromState: &linear.IssueHistoryIssueHistoryIssueHistoryConnectionNodesIssueHistoryFromStateWorkflowState{
			Id:   "22222222-3333-4444-5555-666666666666",
			Name: "Todo",
		},
		ToState: &linear.IssueHistoryIssueHistoryIssueHistoryConnectionNodesIssueHistoryToStateWorkflowState{
			Id:   "33333333-4444-5555-6666-777777777777",
			Name: "In Progress",
		},
	}

	got := MapHistoryEntry(h)

	assertField(t, got, "id", "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	assertField(t, got, "createdAt", "2025-04-01T10:00:00.000Z")

	actor := got["actor"].(map[string]any)
	if actor["name"] != "Ada Lovelace" {
		t.Errorf("actor.name = %v", actor["name"])
	}

	from := got["fromState"].(map[string]any)
	if from["name"] != "Todo" {
		t.Errorf("fromState.name = %v", from["name"])
	}

	to := got["toState"].(map[string]any)
	if to["name"] != "In Progress" {
		t.Errorf("toState.name = %v", to["name"])
	}
}

func TestMapHistoryEntry_LabelChange(t *testing.T) {
	h := HistoryNode{
		Id:        "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		CreatedAt: "2025-04-02T08:00:00.000Z",
		AddedLabels: []linear.IssueHistoryIssueHistoryIssueHistoryConnectionNodesIssueHistoryAddedLabelsIssueLabel{
			{Id: "11111111-2222-3333-4444-555555555555", Name: "bug"},
		},
		RemovedLabels: []linear.IssueHistoryIssueHistoryIssueHistoryConnectionNodesIssueHistoryRemovedLabelsIssueLabel{
			{Id: "22222222-3333-4444-5555-666666666666", Name: "triage"},
		},
	}

	got := MapHistoryEntry(h)

	added := got["addedLabels"].([]map[string]any)
	if len(added) != 1 || added[0]["name"] != "bug" {
		t.Errorf("addedLabels = %v", added)
	}

	removed := got["removedLabels"].([]map[string]any)
	if len(removed) != 1 || removed[0]["name"] != "triage" {
		t.Errorf("removedLabels = %v", removed)
	}
}

func TestMapHistoryEntry_PriorityChange(t *testing.T) {
	h := HistoryNode{
		Id:           "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		CreatedAt:    "2025-04-03T12:00:00.000Z",
		FromPriority: ptr.To(3.0),
		ToPriority:   ptr.To(1.0),
	}

	got := MapHistoryEntry(h)

	if *got["fromPriority"].(*float64) != 3.0 {
		t.Errorf("fromPriority = %v", got["fromPriority"])
	}
	if *got["toPriority"].(*float64) != 1.0 {
		t.Errorf("toPriority = %v", got["toPriority"])
	}
}

func TestMapHistoryEntry_MinimalEntry(t *testing.T) {
	h := HistoryNode{
		Id:        "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		CreatedAt: "2025-04-04T06:00:00.000Z",
	}

	got := MapHistoryEntry(h)

	if _, ok := got["actor"]; ok {
		t.Error("actor should be absent when nil")
	}
	if _, ok := got["fromState"]; ok {
		t.Error("fromState should be absent when nil")
	}

	added := got["addedLabels"].([]map[string]any)
	if len(added) != 0 {
		t.Error("addedLabels should be empty")
	}
}
