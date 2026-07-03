package resolvers

import (
	"strings"
	"testing"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/testutil"
)

func TestResolveDocument_DirectHit(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("DocumentGet", map[string]any{
		"document": map[string]any{"id": "d-uuid", "slugId": "doc-1", "title": "Spec"},
	})

	got, err := ResolveDocument(mock.Client(), "d-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "d-uuid" || got.SlugId != "doc-1" || got.Title != "Spec" {
		t.Errorf("resolved = %+v", got)
	}
}

func TestResolveDocument_SlugFallback(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	// No DocumentGet handler: the direct lookup errors, forcing the list fallback.
	mock.Handle("DocumentList", map[string]any{
		"documents": map[string]any{
			"nodes": []map[string]any{{"id": "d-uuid", "slugId": "doc-1", "title": "Spec"}},
		},
	})

	got, err := ResolveDocument(mock.Client(), "doc-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "d-uuid" {
		t.Errorf("resolved = %+v", got)
	}
}

func TestResolveDocument_NotFound(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("DocumentList", map[string]any{
		"documents": map[string]any{"nodes": []map[string]any{}},
	})

	_, err := ResolveDocument(mock.Client(), "nope")
	if err == nil {
		t.Fatal("expected error for unknown document")
	}
	if !strings.Contains(err.Error(), "document not found") {
		t.Errorf("unexpected message: %v", err)
	}
	var apiErr *apierrors.APIError
	if !apierrors.As(err, &apiErr) || apiErr.FixableBy != apierrors.FixableByAgent {
		t.Errorf("expected FixableByAgent APIError, got %#v", err)
	}
}
