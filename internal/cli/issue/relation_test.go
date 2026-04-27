package issue

import (
	"testing"

	"github.com/shhac/lin/internal/linear"
)

func TestResolveRelationType(t *testing.T) {
	cases := []struct {
		input  string
		want   linear.IssueRelationType
		wantOK bool
	}{
		{"blocks", linear.IssueRelationTypeBlocks, true},
		{"BLOCKS", linear.IssueRelationTypeBlocks, true},
		{"Duplicate", linear.IssueRelationTypeDuplicate, true},
		{"related", linear.IssueRelationTypeRelated, true},
		{"unknown", "", false},
		{"", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, ok := resolveRelationType(tc.input)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && got != tc.want {
				t.Errorf("type = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestInverseRelationLabel(t *testing.T) {
	cases := map[string]string{
		"blocks":    "blocked_by",
		"duplicate": "duplicate",
		"related":   "related",
		"":          "",
	}
	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			if got := inverseRelationLabel(in); got != want {
				t.Errorf("inverseRelationLabel(%q) = %q, want %q", in, got, want)
			}
		})
	}
}
