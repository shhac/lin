package pretty

import (
	"strings"
	"testing"
)

func TestErrorCardPlain(t *testing.T) {
	got := ErrorCard("ENG-999", "issue not found", "check the identifier", plainOpts(60))
	if !strings.Contains(got, "ENG-999") || !strings.Contains(got, "issue not found") {
		t.Errorf("error card missing id or message: %q", got)
	}
	if !strings.Contains(got, "check the identifier") {
		t.Errorf("error card missing hint: %q", got)
	}
	if strings.Contains(got, "\x1b") {
		t.Errorf("error card emitted ANSI with color off: %q", got)
	}
}

func TestErrorCardNoHint(t *testing.T) {
	got := ErrorCard("ENG-1", "boom", "", plainOpts(60))
	if strings.Contains(got, "\n") {
		t.Errorf("hintless error card should be single line: %q", got)
	}
}
