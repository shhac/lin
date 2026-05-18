package linear

import (
	"encoding/json"
	"testing"
)

// TestInputTypes_NoNullFields guards against a regenerate that drops
// ,omitempty from JSON-bound input fields. Linear's API rejects explicit nulls
// in input types with "Argument Validation Error", so any field left nil must
// be absent from the serialized JSON. See scripts/fix-omitempty.go.
func TestInputTypes_NoNullFields(t *testing.T) {
	s := "hi"
	id := "ABC-1"

	cases := []struct {
		name string
		in   any
	}{
		{"CommentCreateInput", CommentCreateInput{Body: &s, IssueId: &id}},
		{"CommentUpdateInput", CommentUpdateInput{Body: &s}},
		{"IssueCreateInput", IssueCreateInput{Title: &s, TeamId: id}},
		{"IssueUpdateInput", IssueUpdateInput{Title: &s}},
		{"DocumentCreateInput", DocumentCreateInput{Title: s}},
		{"DocumentUpdateInput", DocumentUpdateInput{Title: &s}},
		{"AttachmentCreateInput", AttachmentCreateInput{Title: s, Url: s, IssueId: id}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.in)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var m map[string]any
			if err := json.Unmarshal(b, &m); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			for k, v := range m {
				if v == nil {
					t.Errorf("emits null for %q (JSON: %s)", k, b)
				}
			}
		})
	}
}
