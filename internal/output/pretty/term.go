// Package pretty renders Linear entities as human-readable terminal "cards" —
// the --format pretty alternative to the machine formats (json/yaml/ndjson).
// It is pure rendering: callers fetch and map the entity, then hand the mapped
// map[string]any to a domain renderer. Width and color are resolved once into
// Options; all layout math is done on plain (un-colored) text so ANSI codes
// never corrupt alignment.
package pretty

import (
	"os"
	"strings"

	"golang.org/x/term"
)

// Width bounds. Cards stay readable on very wide terminals (capped) and don't
// collapse on very narrow ones (floored).
const (
	defaultWidth = 80
	minWidth     = 40
	maxWidth     = 100
)

// TerminalWidth returns the usable card width: the real terminal width (capped
// to [minWidth, maxWidth]) when stdout is a terminal, else defaultWidth. A
// flagWidth > 0 always wins, so tests and scripts get deterministic output.
func TerminalWidth(flagWidth int) int {
	if flagWidth > 0 {
		return flagWidth
	}
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return clamp(w, minWidth, maxWidth)
	}
	return defaultWidth
}

// ColorCapable reports whether ANSI styling should be emitted: stdout is a
// terminal, NO_COLOR is unset (https://no-color.org), and TERM isn't "dumb".
func ColorCapable() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
