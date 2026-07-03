package pretty

import (
	"testing"
	"time"
)

func TestRelTime(t *testing.T) {
	now := time.Date(2026, 7, 3, 12, 0, 0, 0, time.UTC)
	opts := Options{Now: now}
	cases := []struct {
		name string
		ts   string
		want string
	}{
		{"just now", now.Add(-30 * time.Second).Format(time.RFC3339), "just now"},
		{"minutes", now.Add(-5 * time.Minute).Format(time.RFC3339), "5 minutes ago"},
		{"one hour", now.Add(-time.Hour).Format(time.RFC3339), "1 hour ago"},
		{"hours", now.Add(-3 * time.Hour).Format(time.RFC3339), "3 hours ago"},
		{"days", now.Add(-2 * 24 * time.Hour).Format(time.RFC3339), "2 days ago"},
		{"months", now.Add(-60 * 24 * time.Hour).Format(time.RFC3339), "2 months ago"},
		{"years", now.Add(-2 * 365 * 24 * time.Hour).Format(time.RFC3339), "2 years ago"},
		{"future clamps to just now", now.Add(time.Hour).Format(time.RFC3339), "just now"},
		{"unparseable", "not-a-time", ""},
		{"empty", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := opts.RelTime(tc.ts); got != tc.want {
				t.Errorf("RelTime(%q) = %q, want %q", tc.ts, got, tc.want)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	ok := []string{
		"2026-07-03T12:00:00Z",
		"2026-07-03T12:00:00.123Z",
		"2026-07-03",
	}
	for _, ts := range ok {
		if _, parsed := parseTime(ts); !parsed {
			t.Errorf("parseTime(%q) should succeed", ts)
		}
	}
	for _, ts := range []string{"", "07/03/2026", "garbage"} {
		if _, parsed := parseTime(ts); parsed {
			t.Errorf("parseTime(%q) should fail", ts)
		}
	}
}

func TestPlural(t *testing.T) {
	cases := []struct {
		n    int
		unit string
		want string
	}{
		{0, "minute", "1 minute ago"}, // n<=0 clamps to 1
		{1, "minute", "1 minute ago"},
		{2, "minute", "2 minutes ago"},
		{5, "hour", "5 hours ago"},
	}
	for _, tc := range cases {
		if got := plural(tc.n, tc.unit); got != tc.want {
			t.Errorf("plural(%d, %q) = %q, want %q", tc.n, tc.unit, got, tc.want)
		}
	}
}

func TestDateOnly(t *testing.T) {
	if got := DateOnly("2026-07-03T12:00:00Z"); got != "2026-07-03" {
		t.Errorf("with time: %q", got)
	}
	if got := DateOnly("2026-07-03"); got != "2026-07-03" {
		t.Errorf("already date-only: %q", got)
	}
}
