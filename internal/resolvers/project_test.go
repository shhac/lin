package resolvers

import (
	"strings"
	"testing"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/testutil"
)

func TestResolveProject_DirectHit(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("ProjectGet", map[string]any{
		"project": map[string]any{"id": "p-uuid", "name": "Roadmap", "slugId": "roa-1"},
	})

	got, err := ResolveProject(mock.Client(), "p-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "p-uuid" || got.Name != "Roadmap" || got.SlugId != "roa-1" {
		t.Errorf("resolved = %+v", got)
	}
}

func TestResolveProject_NameFallback(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	// No ProjectGet handler: the direct lookup errors, forcing the list fallback.
	mock.Handle("ProjectList", map[string]any{
		"projects": map[string]any{
			"nodes": []map[string]any{{"id": "p-uuid", "name": "Roadmap", "slugId": "roa-1"}},
		},
	})

	got, err := ResolveProject(mock.Client(), "Roadmap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "p-uuid" {
		t.Errorf("resolved = %+v", got)
	}
}

func TestResolveProject_NotFound(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("ProjectList", map[string]any{
		"projects": map[string]any{"nodes": []map[string]any{}},
	})

	_, err := ResolveProject(mock.Client(), "nope")
	if err == nil {
		t.Fatal("expected error for unknown project")
	}
	if !strings.Contains(err.Error(), "project not found") {
		t.Errorf("unexpected message: %v", err)
	}
	var apiErr *apierrors.APIError
	if !apierrors.As(err, &apiErr) || apiErr.FixableBy != apierrors.FixableByAgent {
		t.Errorf("expected FixableByAgent APIError, got %#v", err)
	}
}
