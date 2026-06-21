package initiative

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output/pretty"
)

func sampleInitiative() map[string]any {
	return map[string]any{
		"id":          "i1",
		"slugId":      "q3-reliability",
		"url":         "https://linear.app/acme/initiative/q3-reliability",
		"name":        "Q3 Reliability Push",
		"status":      linear.InitiativeStatus("Active"),
		"owner":       map[string]any{"id": "u1", "name": "Alex Rivera"},
		"creator":     map[string]any{"id": "u2", "name": "Sam Lee"},
		"health":      linear.InitiativeUpdateHealthType("onTrack"),
		"description": "Cut customer-facing incidents by hardening core flows.",
		"targetDate":  "2026-09-30",
		"startedAt":   "2026-06-01T00:00:00Z",
		"createdAt":   "2026-05-15T00:00:00Z",
		"updatedAt":   "2026-06-19T12:00:00Z",
	}
}

func testOpts(width int) pretty.Options {
	return pretty.Options{Width: width, Color: false, Now: time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)}
}

func TestRenderInitiativeCardGolden(t *testing.T) {
	got := renderInitiativeCard(sampleInitiative(), testOpts(74))
	golden := filepath.Join("testdata", "initiative_card.golden")
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
		t.Errorf("initiative card mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	flat := strings.Join(strings.Fields(got), " ")
	for _, sub := range []string{"Q3 Reliability Push", "Status Active", "Health onTrack", "Description"} {
		if !strings.Contains(flat, sub) {
			t.Errorf("initiative card missing %q", sub)
		}
	}
}
