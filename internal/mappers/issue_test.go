package mappers

import "testing"

func TestMapIssueSummary_Full(t *testing.T) {
	input := IssueSummaryInput{
		ID:            "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		Identifier:    "ENG-101",
		Title:         "Fix widget alignment",
		BranchName:    "eng-101-fix-widget",
		Priority:      2,
		PriorityLabel: "High",
		StateName:     "In Progress",
		StateType:     "started",
		AssigneeID:    "11111111-2222-3333-4444-555555555555",
		AssigneeName:  "Ada Lovelace",
		TeamKey:       "ENG",
	}
	got := MapIssueSummary(input)

	assertField(t, got, "id", "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	assertField(t, got, "identifier", "ENG-101")
	assertField(t, got, "title", "Fix widget alignment")
	assertField(t, got, "branchName", "eng-101-fix-widget")
	assertField(t, got, "status", "In Progress")
	assertField(t, got, "statusType", "started")
	assertField(t, got, "team", "ENG")
	assertField(t, got, "assignee", "Ada Lovelace")
	assertField(t, got, "assigneeId", "11111111-2222-3333-4444-555555555555")

	if got["priority"] != 2 {
		t.Errorf("priority = %v, want 2", got["priority"])
	}
}

func TestMapIssueSummary_NoAssignee(t *testing.T) {
	input := IssueSummaryInput{
		ID:            "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		Identifier:    "ENG-102",
		Title:         "Unassigned task",
		Priority:      0,
		PriorityLabel: "None",
		StateName:     "Backlog",
		StateType:     "backlog",
		TeamKey:       "ENG",
	}
	got := MapIssueSummary(input)

	if _, ok := got["assignee"]; ok {
		t.Error("assignee should be absent when empty")
	}
	if _, ok := got["assigneeId"]; ok {
		t.Error("assigneeId should be absent when empty")
	}
}

func assertField(t *testing.T, m map[string]any, key, want string) {
	t.Helper()
	got, ok := m[key]
	if !ok {
		t.Errorf("missing key %q", key)
		return
	}
	if got != want {
		t.Errorf("%s = %v, want %v", key, got, want)
	}
}
