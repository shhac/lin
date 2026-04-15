package resolvers

import (
	"strings"
	"testing"

	"github.com/shhac/lin/internal/testutil"
)

func TestResolveLabels_SingleByName(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("LabelList", map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Bug", "color": "#ff0000"},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "Feature", "color": "#00ff00"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	ids, err := ResolveLabels(mock.Client(), "Bug", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 || ids[0] != "aaaaaaaa-1111-2222-3333-444444444444" {
		t.Errorf("ids = %v", ids)
	}
}

func TestResolveLabels_ByID(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("LabelList", map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Bug", "color": "#ff0000"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	ids, err := ResolveLabels(mock.Client(), "aaaaaaaa-1111-2222-3333-444444444444", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 || ids[0] != "aaaaaaaa-1111-2222-3333-444444444444" {
		t.Errorf("ids = %v", ids)
	}
}

func TestResolveLabels_CommaSeparated(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("LabelList", map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Bug", "color": "#ff0000"},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "Feature", "color": "#00ff00"},
				{"id": "cccccccc-1111-2222-3333-444444444444", "name": "Enhancement", "color": "#0000ff"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	ids, err := ResolveLabels(mock.Client(), "Bug, Feature", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 IDs, got %d", len(ids))
	}
	if ids[0] != "aaaaaaaa-1111-2222-3333-444444444444" {
		t.Errorf("first ID = %q", ids[0])
	}
	if ids[1] != "bbbbbbbb-1111-2222-3333-444444444444" {
		t.Errorf("second ID = %q", ids[1])
	}
}

func TestResolveLabels_NotFound(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("LabelList", map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Bug", "color": "#ff0000"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	_, err := ResolveLabels(mock.Client(), "NonExistent", "")
	if err == nil {
		t.Fatal("expected error for not found label")
	}
	if !strings.Contains(err.Error(), "Label not found") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "Bug") {
		t.Error("error should list available labels")
	}
}

func TestResolveLabels_Ambiguous(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("LabelList", map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Bug", "color": "#ff0000"},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "Bug", "color": "#ff0000"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	_, err := ResolveLabels(mock.Client(), "Bug", "")
	if err == nil {
		t.Fatal("expected error for ambiguous label")
	}
	if !strings.Contains(err.Error(), "Ambiguous label") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "Tip: use --team") {
		t.Error("error should include team hint when not team-scoped")
	}
}

func TestResolveLabels_WithTeamScope(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	teamID := "aaaaaaaa-1111-2222-3333-444444444444"
	mock.Handle("TeamLabels", map[string]any{
		"team": map[string]any{
			"labels": map[string]any{
				"nodes": []map[string]any{
					{"id": "11111111-aaaa-bbbb-cccc-dddddddddddd", "name": "Team Bug", "color": "#ff0000"},
				},
				"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
			},
		},
	})
	mock.Handle("LabelList", map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "22222222-aaaa-bbbb-cccc-dddddddddddd", "name": "Global Label", "color": "#0000ff"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	ids, err := ResolveLabels(mock.Client(), "Team Bug", teamID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 || ids[0] != "11111111-aaaa-bbbb-cccc-dddddddddddd" {
		t.Errorf("ids = %v", ids)
	}
}

func TestResolveLabels_CaseInsensitive(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("LabelList", map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Bug", "color": "#ff0000"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	ids, err := ResolveLabels(mock.Client(), "bug", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 {
		t.Errorf("expected 1 ID, got %d", len(ids))
	}
}
