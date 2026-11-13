package document

import (
	"context"
	"reflect"
	"strings"
	"testing"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
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
	mock.Handle("ProjectGet", map[string]any{"project": map[string]any{"id": "p-uuid", "name": "Roadmap", "slugId": "roa-1"}})
	mock.Handle("TeamGet", map[string]any{"team": map[string]any{"id": "t-uuid", "key": "ENG", "name": "Engineering"}})
	mock.Handle("InitiativeGet", map[string]any{"initiative": map[string]any{"id": "i-uuid", "name": "Growth", "slugId": "gro-1"}})
	mock.Handle("DocumentUpdate", map[string]any{"documentUpdate": map[string]any{"success": true, "document": payloadDoc}})
	return mock
}

var (
	projectParent    = map[string]any{"id": "p-uuid", "name": "Roadmap", "slugId": "roa-1"}
	teamParent       = map[string]any{"id": "t-uuid", "key": "ENG", "name": "Engineering"}
	initiativeParent = map[string]any{"id": "i-uuid", "name": "Growth"}
)

func TestMoveDocument_EchoesTheNewParent(t *testing.T) {
	cases := []struct {
		kind, target, inputField string
		parent                   map[string]any
		wantID                   string
	}{
		{"project", "p-uuid", "projectId", projectParent, "p-uuid"},
		{"team", "ENG", "teamId", teamParent, "t-uuid"},
		{"initiative", "i-uuid", "initiativeId", initiativeParent, "i-uuid"},
	}
	for _, tc := range cases {
		t.Run(tc.kind, func(t *testing.T) {
			mock := moveMock(t, map[string]any{"id": "d-uuid", tc.kind: tc.parent})
			defer mock.Close()

			got, err := moveDocument(context.Background(), mock.Client(), moveByKind(t, tc.kind), "d-uuid", tc.target)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := map[string]any{"updated": true, "id": "d-uuid", tc.kind: tc.parent}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("result = %#v, want %#v", got, want)
			}
			if input := mock.Variables("DocumentUpdate")["input"]; !reflect.DeepEqual(input, map[string]any{tc.inputField: tc.wantID}) {
				t.Errorf("input = %#v", input)
			}
		})
	}
}

// Linear acknowledges a parent change it doesn't apply (e.g. a null parent),
// so success alone can't be trusted.
func TestMoveDocument_ParentThatDidNotLandIsAnError(t *testing.T) {
	cases := []struct {
		name string
		doc  map[string]any
	}{
		{"old parent kept", map[string]any{"id": "d-uuid", "project": projectParent}},
		{"same kind, other id", map[string]any{"id": "d-uuid", "team": map[string]any{"id": "t-other", "key": "OPS", "name": "Ops"}}},
		{"new parent beside the old one", map[string]any{"id": "d-uuid", "team": teamParent, "project": projectParent}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock := moveMock(t, tc.doc)
			defer mock.Close()

			_, err := moveDocument(context.Background(), mock.Client(), moveByKind(t, "team"), "d-uuid", "ENG")
			var apiErr *apierrors.APIError
			if !apierrors.As(err, &apiErr) || apiErr.FixableBy != apierrors.FixableByHuman || !strings.Contains(apiErr.Error(), "not the requested team") {
				t.Errorf("want a human-fixable parent-mismatch error, got %#v", err)
			}
		})
	}
}

// soleParent looks the move's kind up in the mapper's output, so a kind the
// mapper never emits would fail every move as "not the requested ...".
func TestParentMoveKindsMatchMapperOutput(t *testing.T) {
	populated := linear.DocumentUpdateDocumentUpdateDocumentPayloadDocument{
		Project:    &linear.DocumentUpdateDocumentUpdateDocumentPayloadDocumentProject{Id: "p"},
		Team:       &linear.DocumentUpdateDocumentUpdateDocumentPayloadDocumentTeam{Id: "t"},
		Initiative: &linear.DocumentUpdateDocumentUpdateDocumentPayloadDocumentInitiative{Id: "i"},
		Issue:      &linear.DocumentUpdateDocumentUpdateDocumentPayloadDocumentIssue{Id: "x"},
	}
	parents := mappers.MapDocumentParents(populated)
	for _, m := range parentMoves() {
		if _, ok := parents[m.kind]; !ok {
			t.Errorf("parent move %q is not a MapDocumentParents key", m.kind)
		}
	}
	if len(mappers.MapDocumentParents(linear.DocumentUpdateDocumentUpdateDocumentPayloadDocument{})) != 0 {
		t.Error("a document with no parents set must map to no parents")
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
