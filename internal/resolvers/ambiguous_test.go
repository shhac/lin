package resolvers

import (
	"strings"
	"testing"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/testutil"
)

func projectNode(id, name, slug string) map[string]any {
	return map[string]any{"id": id, "name": name, "slugId": slug, "state": "started", "lead": map[string]any{"name": "Test Lead"}}
}

func TestResolveProject_DuplicateNamesAreAmbiguous(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("ProjectList", map[string]any{
		"projects": map[string]any{"nodes": []map[string]any{
			projectNode("p-1", "Roadmap", "aaa111"),
			projectNode("p-2", "roadmap", "bbb222"),
		}},
	})

	_, err := ResolveProject(mock.Client(), "Roadmap")
	assertAmbiguous(t, err, `ambiguous project: "Roadmap" matches 2 projects`, "aaa111", "bbb222", "use the slug ID or UUID")
	assertFetchedMoreThanOne(t, mock, "ProjectList")
}

func TestResolveProject_SlugBeatsNameMatch(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("ProjectList", map[string]any{
		"projects": map[string]any{"nodes": []map[string]any{
			projectNode("p-1", "bbb222", "aaa111"),
			projectNode("p-2", "Roadmap", "bbb222"),
		}},
	})

	got, err := ResolveProject(mock.Client(), "bbb222")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "p-2" {
		t.Errorf("resolved %+v, want the slug match p-2", got)
	}
}

func TestResolveCustomer_DuplicateNamesAreAmbiguous(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	customer := func(id, slug string) map[string]any {
		return map[string]any{
			"id": id, "name": "Acme", "slugId": slug, "url": "https://example.test/" + slug,
			"approximateNeedCount": 0, "status": map[string]any{"id": "s1", "name": "Active"},
		}
	}
	mock.Handle("CustomerList", map[string]any{
		"customers": map[string]any{"nodes": []map[string]any{customer("c-1", "acme-1"), customer("c-2", "acme-2")}},
	})

	_, err := ResolveCustomer(mock.Client(), "acme")
	assertAmbiguous(t, err, `ambiguous customer: "acme" matches 2 customers`, "acme-1", "acme-2")
	assertFetchedMoreThanOne(t, mock, "CustomerList")
}

// The mock returns every node whatever "first" asks for, so the limit is
// asserted directly: against Linear, first: 1 would hide the second match.
func assertFetchedMoreThanOne(t *testing.T, mock *testutil.MockLinear, op string) {
	t.Helper()
	first, _ := mock.Variables(op)["first"].(float64)
	if first <= 1 {
		t.Errorf("%s requested first: %v; it must fetch more than one to detect a duplicate", op, first)
	}
}

func TestResolveInitiative_DuplicateNamesAreAmbiguous(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("InitiativeList", map[string]any{
		"initiatives": map[string]any{
			"nodes": []map[string]any{
				{"id": "i-1", "name": "Growth", "slugId": "gro-1"},
				{"id": "i-2", "name": "Growth", "slugId": "gro-2"},
			},
			"pageInfo": map[string]any{"hasNextPage": false},
		},
	})

	_, err := ResolveInitiative(mock.Client(), "growth")
	assertAmbiguous(t, err, `ambiguous initiative: "growth" matches 2 initiatives`, "gro-1", "gro-2")
}

func assertAmbiguous(t *testing.T, err error, wants ...string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an ambiguity error")
	}
	for _, want := range wants {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q missing %q", err.Error(), want)
		}
	}
	var apiErr *apierrors.APIError
	if !apierrors.As(err, &apiErr) || apiErr.FixableBy != apierrors.FixableByAgent {
		t.Errorf("expected FixableByAgent APIError, got %#v", err)
	}
}
