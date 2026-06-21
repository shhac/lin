// Package pretty renders Linear entities as human-readable terminal "cards" —
// the --format pretty alternative to the machine formats (json/yaml/ndjson).
// It is pure rendering: callers fetch and map the entity, then hand the mapped
// map[string]any to a domain renderer. Width and color are resolved once into
// Options; all layout math is done on plain (un-colored) text so ANSI codes
// never corrupt alignment.
package pretty

import (
	"os"

	"golang.org/x/term"

	output "github.com/shhac/lib-agent-output"
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

// ColorCapable reports whether ANSI styling should be emitted for the pretty
// card. It defers to the family color decision (the global --color flag, plus
// NO_COLOR/TERM=dumb and whether stdout is a terminal) so every surface — the
// JSON formats and this human renderer — shares one policy.
func ColorCapable() bool {
	return output.Enabled(os.Stdout)
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
