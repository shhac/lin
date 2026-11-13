package initiative

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/shhac/lin/internal/testutil"
)

func ownerMock(t *testing.T, initiative map[string]any) *testutil.MockLinear {
	t.Helper()
	mock := testutil.NewMockLinear()
	mock.Handle("InitiativeGet", map[string]any{"initiative": map[string]any{"id": "i-uuid", "name": "Growth", "slugId": "gro-1"}})
	mock.Handle("UserList", map[string]any{"users": map[string]any{
		"nodes":    []map[string]any{{"id": "u-1", "name": "Test User", "email": "test@example.test", "displayName": "test"}},
		"pageInfo": map[string]any{"hasNextPage": false},
	}})
	mock.Handle("InitiativeUpdate", map[string]any{"initiativeUpdate": map[string]any{"success": true, "initiative": initiative}})
	return mock
}

func TestSetOwner_EchoesTheOwnerLinearReports(t *testing.T) {
	mock := ownerMock(t, map[string]any{"id": "i-uuid", "owner": map[string]any{"id": "u-1", "name": "Test User"}})
	defer mock.Close()

	got, err := setOwner(context.Background(), mock.Client(), "i-uuid", "test@example.test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]any{"updated": true, "id": "i-uuid", "owner": map[string]any{"id": "u-1", "name": "Test User"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestClearOwner_SendsNullAndReportsCleared(t *testing.T) {
	mock := ownerMock(t, map[string]any{"id": "i-uuid"})
	defer mock.Close()

	got, err := clearOwner(context.Background(), mock.Client(), "i-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := map[string]any{"updated": true, "id": "i-uuid", "cleared": "owner"}; !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
	input, _ := mock.Variables("InitiativeUpdate")["input"].(map[string]any)
	if v, present := input["ownerId"]; !present || v != nil || len(input) != 1 {
		t.Errorf("input = %#v, want only ownerId: null", input)
	}
}

func TestClearOwner_ThatDidNotStickIsAnError(t *testing.T) {
	mock := ownerMock(t, map[string]any{"id": "i-uuid", "owner": map[string]any{"id": "u-1", "name": "Test User"}})
	defer mock.Close()

	_, err := clearOwner(context.Background(), mock.Client(), "i-uuid")
	if err == nil || !strings.Contains(err.Error(), "did not remove the owner") {
		t.Errorf("want a not-removed error, got %v", err)
	}
}
