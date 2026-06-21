package project

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

func sp(s string) *string { return &s }

func sampleProject() map[string]any {
	return map[string]any{
		"id":          "p1",
		"slugId":      "checkout-reliability",
		"url":         "https://linear.app/acme/project/checkout-reliability",
		"name":        "Checkout Reliability",
		"description": "Reduce checkout flakiness and improve conversion.",
		"content":     sp("Harden the retry path and add load tests across the checkout flow."),
		"status":      "started",
		"progress":    0.75,
		"startDate":   sp("2026-06-01"),
		"targetDate":  sp("2026-09-30"),
		"lead":        map[string]any{"id": "u1", "name": "Alex Rivera"},
		"labels": []map[string]any{
			{"id": "l1", "name": "infra"},
			{"id": "l2", "name": "reliability"},
		},
		"milestones": []map[string]any{
			{"id": "m1", "name": "Alpha", "targetDate": sp("2026-07-15")},
			{"id": "m2", "name": "Beta", "targetDate": sp("2026-08-20")},
		},
	}
}

func testOpts(width int) pretty.Options {
	return pretty.Options{Width: width, Color: false, Now: time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)}
}

func TestRenderProjectCardGolden(t *testing.T) {
	got := renderProjectCard(sampleProject(), testOpts(74))
	golden := filepath.Join("testdata", "project_card.golden")
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
		t.Errorf("project card mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	flat := strings.Join(strings.Fields(got), " ")
	for _, sub := range []string{"Checkout Reliability", "75% complete", "Status Started", "Milestones", "Alpha"} {
		if !strings.Contains(flat, sub) {
			t.Errorf("project card missing %q", sub)
		}
	}
}

// TestRenderProjectCardColorParity verifies the colored render carries ANSI and
// that stripping it reproduces the plain layout exactly — i.e. coloring a grid
// value (status) doesn't shift column alignment.
func TestRenderProjectCardColorParity(t *testing.T) {
	plain := renderProjectCard(sampleProject(), testOpts(74))
	opts := testOpts(74)
	opts.Color = true
	colored := renderProjectCard(sampleProject(), opts)

	if !strings.Contains(colored, "\x1b[33m") { // yellow "Started" (started state)
		t.Error("expected colored status in project card")
	}
	if stripANSI(colored) != plain {
		t.Errorf("colored visible text differs from plain:\n--- stripped ---\n%s\n--- plain ---\n%s", stripANSI(colored), plain)
	}
}
