package resolvers

import (
	"strings"
	"testing"

	"github.com/shhac/lin/internal/testutil"
)

func TestResolveProjectLabels_SingleByName(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("ProjectLabelList", map[string]any{
		"projectLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Discovery", "color": "#ff0000"},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "Roadmap", "color": "#00ff00"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	ids, err := ResolveProjectLabels(mock.Client(), "Discovery")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 || ids[0] != "aaaaaaaa-1111-2222-3333-444444444444" {
		t.Errorf("ids = %v", ids)
	}
}

func TestResolveProjectLabels_ByID(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("ProjectLabelList", map[string]any{
		"projectLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Discovery", "color": "#ff0000"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	ids, err := ResolveProjectLabels(mock.Client(), "aaaaaaaa-1111-2222-3333-444444444444")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 || ids[0] != "aaaaaaaa-1111-2222-3333-444444444444" {
		t.Errorf("ids = %v", ids)
	}
}

func TestResolveProjectLabels_CommaSeparated(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("ProjectLabelList", map[string]any{
		"projectLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Discovery", "color": "#ff0000"},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "Roadmap", "color": "#00ff00"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	ids, err := ResolveProjectLabels(mock.Client(), "Discovery, Roadmap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 IDs, got %d", len(ids))
	}
}

func TestResolveProjectLabels_NotFound(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("ProjectLabelList", map[string]any{
		"projectLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Discovery", "color": "#ff0000"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	_, err := ResolveProjectLabels(mock.Client(), "NonExistent")
	if err == nil {
		t.Fatal("expected error for not found label")
	}
	if !strings.Contains(err.Error(), "project label not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestResolveProjectLabels_Ambiguous(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("ProjectLabelList", map[string]any{
		"projectLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Discovery", "color": "#ff0000"},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "Discovery", "color": "#00ff00"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	_, err := ResolveProjectLabels(mock.Client(), "Discovery")
	if err == nil {
		t.Fatal("expected error for ambiguous label")
	}
	if !strings.Contains(err.Error(), "ambiguous project label") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestResolveProjectLabels_CaseInsensitive(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("ProjectLabelList", map[string]any{
		"projectLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Discovery", "color": "#ff0000"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	ids, err := ResolveProjectLabels(mock.Client(), "discovery")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 {
		t.Errorf("expected 1 ID, got %d", len(ids))
	}
}
