package customer

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/shhac/lin/internal/output/pretty"
)

var ansiRE = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string { return ansiRE.ReplaceAllString(s, "") }

func sampleCustomer() map[string]any {
	return map[string]any{
		"id":                   "c1",
		"name":                 "Acme Corp",
		"slugId":               "acme-corp",
		"url":                  "https://linear.app/acme/customer/acme-corp",
		"approximateNeedCount": 12,
		"status":               map[string]any{"id": "s1", "name": "Active"},
		"domains":              []any{"acme.com", "acme.io"},
		"externalIds":          []any{"crm-4821"},
		"createdAt":            "2026-01-10T00:00:00Z",
		"updatedAt":            "2026-06-16T12:00:00Z",
		"owner":                map[string]any{"id": "u1", "name": "Alex Rivera"},
		"tier":                 map[string]any{"id": "t1", "displayName": "Enterprise"},
		"revenue":              50000.0,
		"size":                 200.0,
	}
}

func testOpts(width int) pretty.Options {
	return pretty.Options{Width: width, Color: false, Now: time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)}
}

func TestRenderCustomerCardGolden(t *testing.T) {
	got := renderCustomerCard(sampleCustomer(), testOpts(74))
	golden := filepath.Join("testdata", "customer_card.golden")
	if os.Getenv("UPDATE_GOLDEN") != "" {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read golden (UPDATE_GOLDEN=1 to create): %v", err)
	}
	if got != string(want) {
		t.Errorf("customer card mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	flat := strings.Join(strings.Fields(got), " ")
	for _, sub := range []string{"Acme Corp", "Tier Enterprise", "Domains", "acme.com · acme.io"} {
		if !strings.Contains(flat, sub) {
			t.Errorf("customer card missing %q", sub)
		}
	}
}

func TestRenderCustomerCardColorParity(t *testing.T) {
	plain := renderCustomerCard(sampleCustomer(), testOpts(74))
	opts := testOpts(74)
	opts.Color = true
	colored := renderCustomerCard(sampleCustomer(), opts)

	if !strings.Contains(colored, "\x1b[32m") { // green "Active" (heuristic)
		t.Error("expected colored status in customer card")
	}
	if stripANSI(colored) != plain {
		t.Errorf("colored visible text differs from plain:\n--- stripped ---\n%s\n--- plain ---\n%s", stripANSI(colored), plain)
	}
}

// renderCustomerCard reads the raw mapper map (no JSON round-trip), so domains
// and externalIds arrive as the mapper's native []string. This guards against
// joinAny only handling the post-JSON []any shape.
func TestRenderCustomerCard_NativeStringSlices(t *testing.T) {
	d := map[string]any{
		"name":        "Acme Corp",
		"url":         "https://linear.app/acme/customer/acme",
		"status":      map[string]any{"id": "s1", "name": "Active"},
		"domains":     []string{"acme.example", "acme.test"},
		"externalIds": []string{"sf-123"},
	}
	got := renderCustomerCard(d, testOpts(74))
	flat := strings.Join(strings.Fields(got), " ")
	for _, sub := range []string{"Acme Corp", "Domains acme.example · acme.test", "External sf-123"} {
		if !strings.Contains(flat, sub) {
			t.Errorf("customer card missing %q in:\n%s", sub, got)
		}
	}
}

func TestJoinAny(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want string
	}{
		{"native []string", []string{"a", "b"}, "a · b"},
		{"json []any", []any{"a", "b"}, "a · b"},
		{"empty []string", []string{}, ""},
		{"non-slice", "nope", ""},
		{"nil", nil, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := joinAny(tc.in); got != tc.want {
				t.Errorf("joinAny(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
