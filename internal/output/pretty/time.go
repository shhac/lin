package pretty

import (
	"fmt"
	"strings"
	"time"
)

// RelTime renders an ISO-8601 timestamp as a coarse "x ago" relative to now.
// Returns "" for an unparseable/empty stamp so callers can omit the field.
func (o Options) RelTime(ts string) string {
	t, ok := parseTime(ts)
	if !ok {
		return ""
	}
	d := o.now().Sub(t)
	if d < 0 {
		d = 0
	}
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return plural(int(d.Minutes()), "minute")
	case d < 24*time.Hour:
		return plural(int(d.Hours()), "hour")
	case d < 30*24*time.Hour:
		return plural(int(d.Hours()/24), "day")
	case d < 365*24*time.Hour:
		return plural(int(d.Hours()/(24*30)), "month")
	default:
		return plural(int(d.Hours()/(24*365)), "year")
	}
}

// DateOnly returns the calendar-date portion of a timestamp (YYYY-MM-DD),
// passing through a value that is already date-only.
func DateOnly(ts string) string {
	if i := strings.IndexByte(ts, 'T'); i >= 0 {
		return ts[:i]
	}
	return ts
}

func parseTime(ts string) (time.Time, bool) {
	if ts == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, ts); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func plural(n int, unit string) string {
	if n <= 0 {
		n = 1
	}
	if n == 1 {
		return fmt.Sprintf("1 %s ago", unit)
	}
	return fmt.Sprintf("%d %ss ago", n, unit)
}
