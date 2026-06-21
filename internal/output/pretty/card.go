package pretty

import (
	"strings"
	"unicode/utf8"
)

// Box-drawing glyphs used across cards.
const (
	ruleChar = "─"
	quoteBar = "▎" // U+258E, left-bar blockquote
)

// Builder accumulates the lines of a single card. All width math uses visible
// (rune) length on plain text; colored leaves are inserted only after padding is
// computed, so ANSI codes never skew alignment.
type Builder struct {
	opts Options
	b    strings.Builder
}

// New starts a card with the given options.
func New(opts Options) *Builder { return &Builder{opts: opts} }

// String returns the assembled card without a trailing newline.
func (c *Builder) String() string {
	return strings.TrimRight(c.b.String(), "\n")
}

func (c *Builder) line(s string) {
	c.b.WriteString(s)
	c.b.WriteByte('\n')
}

// Blank emits an empty line.
func (c *Builder) Blank() { c.b.WriteByte('\n') }

// Rule draws a full-width horizontal rule.
func (c *Builder) Rule() { c.line(strings.Repeat(ruleChar, c.opts.Width)) }

// Separator returns a full-width rule string, for callers placing dividers
// between stacked cards in a multi-get.
func Separator(opts Options) string { return strings.Repeat(ruleChar, opts.Width) }

// Header renders the top line: left content, then right content pushed to the
// right edge. If they'd overlap, they're separated by a single space instead.
// left/right may already contain ANSI codes; plainLeft/plainRight give their
// visible text for the spacing math.
func (c *Builder) Header(left, plainLeft, right, plainRight string) {
	gap := c.opts.Width - runeLen(plainLeft) - runeLen(plainRight)
	if gap < 1 {
		gap = 1
	}
	c.line(left + strings.Repeat(" ", gap) + right)
}

// Title renders a plain, optionally bold, single line.
func (c *Builder) Title(s string) { c.line(s) }

// Line emits a pre-composed line verbatim (the caller owns any styling/width).
func (c *Builder) Line(s string) { c.line(s) }

// Width returns the card width, for callers composing their own lines.
func (c *Builder) Width() int { return c.opts.Width }

// Section draws a titled rule: ─ Title ───────.
func (c *Builder) Section(title string) {
	prefix := ruleChar + " " + title + " "
	fill := c.opts.Width - runeLen(prefix)
	if fill < 0 {
		fill = 0
	}
	c.line(c.opts.Dim(prefix + strings.Repeat(ruleChar, fill)))
}

// Wrapped word-wraps s to the card width and emits each line.
func (c *Builder) Wrapped(s string) {
	for _, ln := range wrap(s, c.opts.Width) {
		c.line(ln)
	}
}

// Grid renders label/value pairs in two columns, row-major, omitting nothing
// (callers filter empties first). Labels are dimmed; the second column starts at
// a fixed offset so values align.
func (c *Builder) Grid(pairs [][2]string) {
	if len(pairs) == 0 {
		return
	}
	labelW := 0
	for _, p := range pairs {
		if n := runeLen(p[0]); n > labelW {
			labelW = n
		}
	}
	colW := c.opts.Width / 2
	cell := func(p [2]string) (text string, plainW int) {
		label := p[0] + strings.Repeat(" ", labelW-runeLen(p[0]))
		plain := label + "  " + p[1]
		return c.opts.Dim(label) + "  " + p[1], runeLen(plain)
	}
	for i := 0; i < len(pairs); i += 2 {
		text, plainW := cell(pairs[i])
		if i+1 < len(pairs) {
			pad := colW - plainW
			if pad < 2 {
				pad = 2
			}
			right, _ := cell(pairs[i+1])
			c.line(text + strings.Repeat(" ", pad) + right)
			continue
		}
		c.line(text)
	}
}

// Field renders a single full-width "Label  value" line (used for things too
// wide for the grid, e.g. labels list).
func (c *Builder) Field(label, value string) {
	c.line(c.opts.Dim(label) + "  " + value)
}

// Blockquote renders a header line followed by body text indented under a
// vertical bar, wrapped to width. Used for comments under --full.
func (c *Builder) Blockquote(header, body string) {
	c.line(header)
	indent := quoteBar + " "
	for _, ln := range wrap(body, c.opts.Width-runeLen(indent)) {
		c.line(c.opts.Dim(quoteBar) + " " + ln)
	}
}

func runeLen(s string) int { return utf8.RuneCountInString(s) }

// Capitalize upper-cases the first rune of s, leaving the rest unchanged. Used
// to present lowercase enum values (e.g. a project state "started") as "Started".
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if r[0] >= 'a' && r[0] <= 'z' {
		r[0] -= 'a' - 'A'
	}
	return string(r)
}

// wrap breaks s into lines no wider than width, splitting on spaces and
// preserving existing newlines (paragraphs). A word longer than width is left
// intact on its own line rather than hard-split.
func wrap(s string, width int) []string {
	if width < 1 {
		width = 1
	}
	var out []string
	for _, para := range strings.Split(s, "\n") {
		words := strings.Fields(para)
		if len(words) == 0 {
			out = append(out, "")
			continue
		}
		line := words[0]
		for _, w := range words[1:] {
			if runeLen(line)+1+runeLen(w) > width {
				out = append(out, line)
				line = w
				continue
			}
			line += " " + w
		}
		out = append(out, line)
	}
	return out
}
