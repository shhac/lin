package pretty

import "strings"

// ErrorCard renders a single failed lookup in a multi-get as a compact,
// human-readable block instead of the machine @unresolved record. id is the
// requested identifier, msg the failure reason, hint an optional next step.
func ErrorCard(id, msg, hint string, opts Options) string {
	var b strings.Builder
	mark := opts.style(ansiRed, "✗")
	b.WriteString(mark + " " + opts.Bold(id))
	if msg != "" {
		b.WriteString(opts.Dim(" — ") + msg)
	}
	if hint != "" {
		b.WriteByte('\n')
		for _, ln := range wrap(hint, opts.Width-2) {
			b.WriteString("  " + opts.Dim(ln) + "\n")
		}
		return strings.TrimRight(b.String(), "\n")
	}
	return b.String()
}
