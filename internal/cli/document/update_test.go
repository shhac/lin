package document

import (
	"context"
	"reflect"
	"strings"
	"testing"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/testutil"
)

func moveByKind(t *testing.T, kind string) parentMove {
	t.Helper()
	for _, m := range parentMoves() {
		if m.kind == kind {
			return m
		}
	}
	t.Fatalf("no parent move %q", kind)
	return parentMove{}
}

func moveMock(t *testing.T, payloadDoc map[string]any) *testutil.MockLinear {
	t.Helper()
	mock := testutil.NewMockLinear()
	mock.Handle("DocumentGet", map[string]any{"document": map[string]any{"id": "d-uuid", "slugId": "doc-1", "title": "Spec"}})
	mock.Handle("TeamGet", map[string]any{"team": map[string]any{"id": "t-uuid", "key": "ENG", "name": "Engineering"}})
	mock.Handle("DocumentUpdate", map[string]any{"documentUpdate": map[string]any{"success": true, "document": payloadDoc}})
	return mock
}

func TestMoveDocument_EchoesTheNewParent(t *testing.T) {
	mock := moveMock(t, map[string]any{"id": "d-uuid", "team": map[string]any{"id": "t-uuid", "key": "ENG", "name": "Engineering"}})
	defer mock.Close()

	got, err := moveDocument(context.Background(), mock.Client(), moveByKind(t, "team"), "d-uuid", "ENG")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]any{"updated": true, "id": "d-uuid", "team": map[string]any{"id": "t-uuid", "key": "ENG", "name": "Engineering"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
	if input := mock.Variables("DocumentUpdate")["input"]; !reflect.DeepEqual(input, map[string]any{"teamId": "t-uuid"}) {
		t.Errorf("input = %#v", input)
	}
}

// Linear acknowledges a parent change it doesn't apply (e.g. a null parent),
// so success alone can't be trusted.
func TestMoveDocument_ParentThatDidNotChangeIsAnError(t *testing.T) {
	mock := moveMock(t, map[string]any{"id": "d-uuid", "project": map[string]any{"id": "p-uuid", "name": "Roadmap", "slugId": "roa-1"}})
	defer mock.Close()

	_, err := moveDocument(context.Background(), mock.Client(), moveByKind(t, "team"), "d-uuid", "ENG")
	var apiErr *apierrors.APIError
	if !apierrors.As(err, &apiErr) || apiErr.FixableBy != apierrors.FixableByHuman || !strings.Contains(apiErr.Error(), "not the requested team") {
		t.Errorf("want a human-fixable parent-mismatch error, got %#v", err)
	}
}

func TestNoParentlessDocumentErr_PointsAtMoving(t *testing.T) {
	err := noParentlessDocumentErr("project")
	var apiErr *apierrors.APIError
	if !apierrors.As(err, &apiErr) {
		t.Fatalf("want APIError, got %#v", err)
	}
	if !strings.Contains(apiErr.Error(), "can't be removed from its project") || !strings.Contains(apiErr.Hint, "lin document update project|team|initiative") {
		t.Errorf("got %#v", apiErr)
	}
}

func TestEveryParentMoveSetsItsOwnInputField(t *testing.T) {
	fields := map[string]string{"project": "ProjectId", "team": "TeamId", "initiative": "InitiativeId"}
	for _, m := range parentMoves() {
		in := m.input("x")
		set := map[string]bool{"ProjectId": in.ProjectId != nil, "TeamId": in.TeamId != nil, "InitiativeId": in.InitiativeId != nil}
		for field, isSet := range set {
			if isSet != (field == fields[m.kind]) {
				t.Errorf("%s move: %s set = %v", m.kind, field, isSet)
			}
		}
	}
}
