package issue

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/shhac/lin/internal/output/pretty"
)

func strptr(s string) *string { return &s }
func fptr(f float64) *float64 { return &f }

// sampleIssue mirrors the native-typed map mappers.MapIssueDetail produces
// (pointers for nullable fields), with fully synthetic placeholder data.
func sampleIssue() map[string]any {
	return map[string]any{
		"id":            "issue-uuid",
		"identifier":    "ENG-123",
		"url":           "https://linear.app/acme/issue/ENG-123",
		"title":         "Fix flaky checkout test",
		"description":   strptr("The checkout test fails intermittently under load.\n\nLikely a timing issue in the retry path."),
		"branchName":    "alex/eng-123-fix-flaky-checkout-test",
		"status":        map[string]any{"id": "s1", "name": "In Progress", "type": "started"},
		"assignee":      map[string]any{"id": "u1", "name": "Alex Rivera"},
		"team":          map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
		"project":       map[string]any{"id": "p1", "name": "Checkout Reliability"},
		"parent":        map[string]any{"id": "pi", "identifier": "ENG-100"},
		"priority":      float64(2),
		"priorityLabel": "High",
		"estimate":      fptr(3),
		"dueDate":       strptr("2026-06-30"),
		"labels": []map[string]any{
			{"id": "l1", "name": "bug"},
			{"id": "l2", "name": "flaky"},
		},
		"attachments": []map[string]any{
			{"title": "Fix #4521 checkout retry", "url": "https://github.com/acme/repo/pull/4521", "sourceType": "github"},
		},
		"commentCount":           2,
		"customerRequestCount":   3,
		"customerImportantCount": 1,
		"createdAt":              "2026-06-18T12:00:00Z",
		"updatedAt":              "2026-06-21T10:00:00Z",
	}
}

func testOpts(width int) pretty.Options {
	return pretty.Options{
		Width: width,
		Color: false,
		Now:   time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC),
	}
}

// TestRenderIssueCardGolden snapshots the full card at a fixed width with color
// off. Regenerate with: UPDATE_GOLDEN=1 go test ./internal/cli/issue/.
func TestRenderIssueCardGolden(t *testing.T) {
	const width = 74
	got := renderIssueCard(sampleIssue(), testOpts(width))

	golden := filepath.Join("testdata", "issue_card.golden")
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
		t.Fatalf("read golden (run with UPDATE_GOLDEN=1 to create): %v", err)
	}
	if got != string(want) {
		t.Errorf("card mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}

	if strings.Contains(got, "\x1b") {
		t.Error("ANSI escape present with color off")
	}
	for _, ln := range strings.Split(got, "\n") {
		if utf8.RuneCountInString(ln) > width {
			t.Errorf("line exceeds width %d: %q (%d)", width, ln, utf8.RuneCountInString(ln))
		}
	}
}

// TestRenderIssueCardMinimal covers an unassigned issue with no project, parent,
// labels, attachments, description, or activity — every optional row omitted.
func TestRenderIssueCardMinimal(t *testing.T) {
	d := map[string]any{
		"identifier":    "ENG-9",
		"url":           "https://linear.app/acme/issue/ENG-9",
		"title":         "Bare issue",
		"branchName":    "alex/eng-9-bare-issue",
		"status":        map[string]any{"name": "Backlog", "type": "backlog"},
		"team":          map[string]any{"key": "ENG", "name": "Engineering"},
		"priority":      float64(0),
		"priorityLabel": "No priority",
		"createdAt":     "2026-06-21T11:00:00Z",
		"updatedAt":     "2026-06-21T11:30:00Z",
	}
	got := renderIssueCard(d, testOpts(74))

	if strings.Contains(got, "Assignee  Unassigned") == false {
		t.Errorf("expected Unassigned assignee row, got:\n%s", got)
	}
	for _, absent := range []string{"Project", "Parent", "Labels", "Description", "Attachments", "Activity", "· High", "pts"} {
		if strings.Contains(got, absent) {
			t.Errorf("expected %q omitted from minimal card, got:\n%s", absent, got)
		}
	}
}

func TestRenderIssueCardColorOn(t *testing.T) {
	opts := testOpts(74)
	opts.Color = true
	got := renderIssueCard(sampleIssue(), opts)
	if !strings.Contains(got, "\x1b[1m") {
		t.Error("expected bold ANSI in colored card")
	}
}
