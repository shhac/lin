package resolvers

import (
	"strings"
	"testing"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/testutil"
)

func TestResolveInitiative_DirectHit(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("InitiativeGet", map[string]any{
		"initiative": map[string]any{"id": "i-uuid", "name": "Q3 Roadmap", "slugId": "q3-1"},
	})

	got, err := ResolveInitiative(mock.Client(), "i-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "i-uuid" || got.Name != "Q3 Roadmap" || got.SlugId != "q3-1" {
		t.Errorf("resolved = %+v", got)
	}
}

func TestResolveInitiative_NameFallbackCaseInsensitive(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	// No InitiativeGet handler: the direct lookup errors, forcing the list
	// fallback, whose name match is case-insensitive.
	mock.Handle("InitiativeList", map[string]any{
		"initiatives": map[string]any{
			"nodes": []map[string]any{
				{"id": "i-uuid", "name": "Q3 Roadmap", "slugId": "q3-1"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	got, err := ResolveInitiative(mock.Client(), "q3 roadmap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "i-uuid" {
		t.Errorf("resolved = %+v", got)
	}
}

func TestResolveInitiative_NotFoundListsChoices(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("InitiativeList", map[string]any{
		"initiatives": map[string]any{
			"nodes": []map[string]any{
				{"id": "i-uuid", "name": "Q3 Roadmap", "slugId": "q3-1"},
				{"id": "j-uuid", "name": "Q4 Roadmap", "slugId": "q4-1"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	_, err := ResolveInitiative(mock.Client(), "nope")
	if err == nil {
		t.Fatal("expected error for unknown initiative")
	}
	if !strings.Contains(err.Error(), "initiative not found") || !strings.Contains(err.Error(), "Q3 Roadmap (q3-1)") {
		t.Errorf("expected not-found message listing choices, got: %v", err)
	}
	var apiErr *apierrors.APIError
	if !apierrors.As(err, &apiErr) || apiErr.FixableBy != apierrors.FixableByAgent {
		t.Errorf("expected FixableByAgent APIError, got %#v", err)
	}
}
