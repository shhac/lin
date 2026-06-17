package resolvers

import (
	"strings"
	"testing"

	"github.com/shhac/lin/internal/testutil"
)

const testCustomerID = "cccccccc-1111-2222-3333-444444444444"

func TestResolveCustomer_UUIDDirectLookup(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("CustomerGet", map[string]any{
		"customer": map[string]any{
			"id":                   testCustomerID,
			"name":                 "Acme Corp",
			"slugId":               "acme-corp",
			"url":                  "https://linear.app/acme/customer/acme-corp",
			"domains":              []string{"acme.example"},
			"externalIds":          []string{},
			"revenue":              nil,
			"size":                 nil,
			"approximateNeedCount": 3,
			"createdAt":            "2026-01-01T00:00:00.000Z",
			"updatedAt":            "2026-01-02T00:00:00.000Z",
			"owner":                nil,
			"status":               map[string]any{"id": "s1", "name": "Active"},
			"tier":                 nil,
		},
	})

	got, err := ResolveCustomer(mock.Client(), testCustomerID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != testCustomerID {
		t.Errorf("ID = %q", got.ID)
	}
	if got.SlugId != "acme-corp" {
		t.Errorf("SlugId = %q", got.SlugId)
	}
}

func TestResolveCustomer_NameFallback(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	// No CustomerGet handler: the direct lookup fails, falling back to CustomerList.
	mock.Handle("CustomerList", map[string]any{
		"customers": map[string]any{
			"nodes": []map[string]any{
				{
					"id":                   testCustomerID,
					"name":                 "Acme Corp",
					"slugId":               "acme-corp",
					"url":                  "https://linear.app/acme/customer/acme-corp",
					"revenue":              nil,
					"approximateNeedCount": 3,
					"owner":                nil,
					"status":               map[string]any{"id": "s1", "name": "Active"},
					"tier":                 nil,
				},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	got, err := ResolveCustomer(mock.Client(), "Acme Corp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Acme Corp" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestResolveCustomer_NotFound(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("CustomerList", map[string]any{
		"customers": map[string]any{
			"nodes":    []map[string]any{},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	_, err := ResolveCustomer(mock.Client(), "Nonexistent")
	if err == nil {
		t.Fatal("expected error for not found")
	}
	if !strings.Contains(err.Error(), "customer not found") {
		t.Errorf("unexpected error: %v", err)
	}
}
