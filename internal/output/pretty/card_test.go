package pretty

import (
	"strings"
	"testing"
)

func plainOpts(width int) Options { return Options{Width: width, Color: false} }

func TestWrapRespectsWidth(t *testing.T) {
	lines := wrap("the quick brown fox jumps over the lazy dog", 15)
	for _, ln := range lines {
		if runeLen(ln) > 15 {
			t.Errorf("line exceeds width 15: %q (%d)", ln, runeLen(ln))
		}
	}
	if strings.Join(lines, " ") != "the quick brown fox jumps over the lazy dog" {
		t.Errorf("wrap altered content: %q", lines)
	}
}

func TestWrapPreservesParagraphs(t *testing.T) {
	lines := wrap("first\n\nthird", 40)
	want := []string{"first", "", "third"}
	if strings.Join(lines, "|") != strings.Join(want, "|") {
		t.Errorf("got %q, want %q", lines, want)
	}
}

func TestWrapLongWordNotSplit(t *testing.T) {
	lines := wrap("supercalifragilistic", 5)
	if len(lines) != 1 || lines[0] != "supercalifragilistic" {
		t.Errorf("long word should stay intact, got %q", lines)
	}
}

func TestHeaderRightAligns(t *testing.T) {
	c := New(plainOpts(40))
	c.Header("LEFT", "LEFT", "RIGHT", "RIGHT")
	line := strings.TrimRight(c.String(), "\n")
	if runeLen(line) != 40 {
		t.Errorf("header width = %d, want 40: %q", runeLen(line), line)
	}
	if !strings.HasPrefix(line, "LEFT") || !strings.HasSuffix(line, "RIGHT") {
		t.Errorf("header not aligned: %q", line)
	}
}

func TestRuleWidth(t *testing.T) {
	c := New(plainOpts(20))
	c.Rule()
	if got := strings.TrimRight(c.String(), "\n"); runeLen(got) != 20 {
		t.Errorf("rule width = %d, want 20", runeLen(got))
	}
}

func TestColorGatingOff(t *testing.T) {
	o := plainOpts(40)
	if got := o.Bold("x"); got != "x" {
		t.Errorf("Bold with color off = %q, want %q", got, "x")
	}
	if strings.Contains(o.StatusStyle("started", "In Progress"), "\x1b") {
		t.Error("StatusStyle emitted ANSI with color off")
	}
}

func TestColorGatingOn(t *testing.T) {
	o := Options{Width: 40, Color: true}
	if !strings.Contains(o.Bold("x"), "\x1b[1m") {
		t.Error("Bold with color on should emit ANSI")
	}
}
