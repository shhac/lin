package pretty

import "time"

// Options controls a render pass. Build it once (from resolved width/color) and
// thread it through every renderer.
type Options struct {
	Width int       // total card width in columns
	Color bool      // emit ANSI styling
	Now   time.Time // reference point for relative times; zero means time.Now()
}

// now returns the reference time, defaulting to the wall clock when unset. Tests
// pin Now for deterministic "x ago" output.
func (o Options) now() time.Time {
	if o.Now.IsZero() {
		return time.Now()
	}
	return o.Now
}

// ANSI SGR codes. Applied only to leaf strings whose width is never measured,
// so layout math stays correct.
const (
	ansiReset   = "\x1b[0m"
	ansiBold    = "\x1b[1m"
	ansiDim     = "\x1b[2m"
	ansiRed     = "\x1b[31m"
	ansiGreen   = "\x1b[32m"
	ansiYellow  = "\x1b[33m"
	ansiBlue    = "\x1b[34m"
	ansiMagenta = "\x1b[35m"
	ansiCyan    = "\x1b[36m"
)

func (o Options) style(code, s string) string {
	if !o.Color || s == "" {
		return s
	}
	return code + s + ansiReset
}

// Exported styling helpers, used by domain renderers in their own packages.

// Bold renders s bold when color is enabled.
func (o Options) Bold(s string) string { return o.style(ansiBold, s) }

// Dim renders s dim/faint when color is enabled.
func (o Options) Dim(s string) string { return o.style(ansiDim, s) }

// Accent renders s in an accent color (links, branch names).
func (o Options) Accent(s string) string { return o.style(ansiCyan, s) }

// StatusStyle colors a workflow-state name by its Linear state type.
func (o Options) StatusStyle(stateType, name string) string {
	switch stateType {
	case "completed":
		return o.style(ansiGreen, name)
	case "started":
		return o.style(ansiYellow, name)
	case "canceled", "duplicate":
		return o.style(ansiDim, name)
	case "triage":
		return o.style(ansiMagenta, name)
	default: // backlog, unstarted
		return o.style(ansiBlue, name)
	}
}

// HealthStyle colors an initiative/project health value (onTrack, atRisk,
// offTrack) green/yellow/red. Unknown values pass through uncolored.
func (o Options) HealthStyle(health string) string {
	switch health {
	case "onTrack":
		return o.style(ansiGreen, health)
	case "atRisk":
		return o.style(ansiYellow, health)
	case "offTrack":
		return o.style(ansiRed, health)
	default:
		return health
	}
}
