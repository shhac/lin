package resolvers

import (
	"strings"
	"testing"

	"github.com/shhac/lin/internal/testutil"
)

func TestResolveWorkflowState_CaseInsensitiveMatch(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("WorkflowStates", map[string]any{
		"workflowStates": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Backlog", "type": "backlog", "color": "#bbb", "position": 0},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "In Progress", "type": "started", "color": "#0f0", "position": 1},
				{"id": "cccccccc-1111-2222-3333-444444444444", "name": "Done", "type": "completed", "color": "#00f", "position": 2},
			},
		},
	})

	teamID := "eeeeeeee-1111-2222-3333-444444444444"

	t.Run("exact case", func(t *testing.T) {
		got, err := ResolveWorkflowState(mock.Client(), "In Progress", teamID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "In Progress" {
			t.Errorf("Name = %q", got.Name)
		}
		if got.ID != "bbbbbbbb-1111-2222-3333-444444444444" {
			t.Errorf("ID = %q", got.ID)
		}
	})

	t.Run("lowercase", func(t *testing.T) {
		got, err := ResolveWorkflowState(mock.Client(), "backlog", teamID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "Backlog" {
			t.Errorf("Name = %q", got.Name)
		}
	})

	t.Run("uppercase", func(t *testing.T) {
		got, err := ResolveWorkflowState(mock.Client(), "DONE", teamID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "Done" {
			t.Errorf("Name = %q", got.Name)
		}
	})
}

func TestResolveWorkflowState_NotFound(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("WorkflowStates", map[string]any{
		"workflowStates": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Backlog", "type": "backlog", "color": "#bbb", "position": 0},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "Done", "type": "completed", "color": "#00f", "position": 1},
			},
		},
	})

	teamID := "eeeeeeee-1111-2222-3333-444444444444"
	_, err := ResolveWorkflowState(mock.Client(), "In Review", teamID)
	if err == nil {
		t.Fatal("expected error for not found status")
	}
	if !strings.Contains(err.Error(), "unknown status") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "Backlog") || !strings.Contains(err.Error(), "Done") {
		t.Error("error should list valid state names")
	}
}
