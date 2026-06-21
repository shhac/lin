package document

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shhac/lin/internal/output/pretty"
)

func sampleDocument() map[string]any {
	return map[string]any{
		"id":        "d1",
		"slugId":    "checkout-runbook",
		"title":     "Checkout Runbook",
		"content":   "Steps to follow when checkout latency spikes during peak load.",
		"url":       "https://linear.app/acme/document/checkout-runbook",
		"icon":      "",
		"color":     "",
		"createdAt": "2026-05-10T00:00:00Z",
		"updatedAt": "2026-06-18T12:00:00Z",
		"project":   map[string]any{"id": "p1", "name": "Checkout Reliability", "slugId": "checkout-reliability"},
		"creator":   map[string]any{"id": "u1", "name": "Alex Rivera"},
		"updatedBy": map[string]any{"id": "u2", "name": "Sam Lee"},
	}
}

func testOpts(width int) pretty.Options {
	return pretty.Options{Width: width, Color: false, Now: time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)}
}

func TestRenderDocumentCardGolden(t *testing.T) {
	got := renderDocumentCard(sampleDocument(), testOpts(74))
	golden := filepath.Join("testdata", "document_card.golden")
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
		t.Errorf("document card mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	flat := strings.Join(strings.Fields(got), " ")
	for _, sub := range []string{"Checkout Runbook", "Project Checkout Reliability", "Updated by Sam Lee", "Content"} {
		if !strings.Contains(flat, sub) {
			t.Errorf("document card missing %q", sub)
		}
	}
}
