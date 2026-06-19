package output

import (
	"testing"

	apierrors "github.com/shhac/lin/internal/errors"
)

// TestParseFormat_Lenient pins the now-shared lenient parser: "ndjson"/"yml"
// aliases and case/whitespace insensitivity are accepted, while a genuinely
// unknown value is still rejected and classified fixable_by:agent.
func TestParseFormat_Lenient(t *testing.T) {
	cases := map[string]Format{
		"json":    FormatJSON,
		"JSON":    FormatJSON,
		"yaml":    FormatYAML,
		"yml":     FormatYAML,
		"YML":     FormatYAML,
		"jsonl":   FormatNDJSON,
		"ndjson":  FormatNDJSON,
		" jsonl ": FormatNDJSON,
	}
	for in, want := range cases {
		got, err := ParseFormat(in)
		if err != nil {
			t.Errorf("ParseFormat(%q) returned error: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("ParseFormat(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseFormat_InvalidIsAgentFixable(t *testing.T) {
	_, err := ParseFormat("bogus")
	if err == nil {
		t.Fatal("ParseFormat(\"bogus\") should return an error")
	}
	var aerr *apierrors.APIError
	if !apierrors.As(err, &aerr) {
		t.Fatalf("error should be an *APIError, got %T", err)
	}
	if aerr.FixableBy != apierrors.FixableByAgent {
		t.Errorf("FixableBy = %q, want %q", aerr.FixableBy, apierrors.FixableByAgent)
	}
	// lin attaches an actionable hint on a bad --format; pin it so the
	// delegation to the shared parser can't silently drop it again.
	if aerr.Hint == "" {
		t.Error("unknown-format error should carry a hint")
	}
}
