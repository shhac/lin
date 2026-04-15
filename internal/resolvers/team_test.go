package resolvers

import (
	"strings"
	"testing"

	"github.com/shhac/lin/internal/testutil"
)

func TestResolveTeam_UUIDDirectLookup(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("TeamGet", map[string]any{
		"team": map[string]any{
			"id":                       "aaaaaaaa-1111-2222-3333-444444444444",
			"name":                     "Engineering",
			"key":                      "ENG",
			"description":              nil,
			"issueEstimationType":      "fibonacci",
			"issueEstimationAllowZero": false,
			"issueEstimationExtended":  false,
			"defaultIssueEstimate":     0,
		},
	})

	got, err := ResolveTeam(mock.Client(), "aaaaaaaa-1111-2222-3333-444444444444")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "aaaaaaaa-1111-2222-3333-444444444444" {
		t.Errorf("ID = %q", got.ID)
	}
	if got.Key != "ENG" {
		t.Errorf("Key = %q", got.Key)
	}
}

func TestResolveTeam_KeyLookup(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	// TeamGet will fail for non-UUID (no handler registered).
	// TeamList returns the matching team.
	mock.Handle("TeamList", map[string]any{
		"teams": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Engineering", "key": "ENG"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	got, err := ResolveTeam(mock.Client(), "ENG")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Key != "ENG" {
		t.Errorf("Key = %q", got.Key)
	}
}

func TestResolveTeam_NotFound(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	// TeamGet will fail with GraphQL error for non-UUID
	mock.Handle("TeamList", map[string]any{
		"teams": map[string]any{
			"nodes":    []map[string]any{},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	_, err := ResolveTeam(mock.Client(), "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for not found")
	}
	if !strings.Contains(err.Error(), "Team not found") {
		t.Errorf("unexpected error: %v", err)
	}
}
