package project

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/shhac/lin/internal/testutil"
)

func leadMock(t *testing.T, payloadProject any) *testutil.MockLinear {
	t.Helper()
	mock := testutil.NewMockLinear()
	mock.Handle("ProjectGet", map[string]any{"project": map[string]any{"id": "p-uuid", "name": "Roadmap", "slugId": "roa-1"}})
	mock.Handle("UserList", map[string]any{"users": map[string]any{
		"nodes":    []map[string]any{{"id": "u-1", "name": "Test User", "email": "test@example.test", "displayName": "test"}},
		"pageInfo": map[string]any{"hasNextPage": false},
	}})
	mock.Handle("ProjectUpdate", map[string]any{"projectUpdate": map[string]any{"success": true, "project": payloadProject}})
	return mock
}

func TestSetLead_EchoesTheLeadLinearReports(t *testing.T) {
	mock := leadMock(t, map[string]any{"id": "p-uuid", "lead": map[string]any{"id": "u-1", "name": "Test User"}})
	defer mock.Close()

	got, err := setLead(context.Background(), mock.Client(), "p-uuid", "test@example.test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]any{"updated": true, "id": "p-uuid", "lead": map[string]any{"id": "u-1", "name": "Test User"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
	if input := mock.Variables("ProjectUpdate")["input"]; !reflect.DeepEqual(input, map[string]any{"leadId": "u-1"}) {
		t.Errorf("input = %#v", input)
	}
}

func TestClearLead_SendsNullAndReportsCleared(t *testing.T) {
	mock := leadMock(t, map[string]any{"id": "p-uuid"})
	defer mock.Close()

	got, err := clearLead(context.Background(), mock.Client(), "p-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := map[string]any{"updated": true, "id": "p-uuid", "cleared": "lead"}; !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
	input, _ := mock.Variables("ProjectUpdate")["input"].(map[string]any)
	if v, present := input["leadId"]; !present || v != nil || len(input) != 1 {
		t.Errorf("input = %#v, want only leadId: null", input)
	}
}

func TestClearLead_ThatDidNotStickIsAnError(t *testing.T) {
	mock := leadMock(t, map[string]any{"id": "p-uuid", "lead": map[string]any{"id": "u-1", "name": "Test User"}})
	defer mock.Close()

	_, err := clearLead(context.Background(), mock.Client(), "p-uuid")
	if err == nil || !strings.Contains(err.Error(), "did not remove the lead") {
		t.Errorf("want a not-removed error, got %v", err)
	}
}

func TestClearLead_WithoutProjectInPayloadDoesNotClaimCleared(t *testing.T) {
	mock := leadMock(t, nil)
	defer mock.Close()

	got, err := clearLead(context.Background(), mock.Client(), "p-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := map[string]any{"updated": true, "id": "p-uuid"}; !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}
