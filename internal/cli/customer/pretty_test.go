package customer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shhac/lin/internal/output/pretty"
)

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
