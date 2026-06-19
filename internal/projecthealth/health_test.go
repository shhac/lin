package projecthealth

import (
	"strings"
	"testing"

	"github.com/shhac/lin/internal/linear"
)

func TestValidate_AcceptedForms(t *testing.T) {
	tests := []struct {
		input string
		want  linear.ProjectUpdateHealthType
	}{
		{"on-track", linear.ProjectUpdateHealthTypeOntrack},
		{"On-Track", linear.ProjectUpdateHealthTypeOntrack},
		{"onTrack", linear.ProjectUpdateHealthTypeOntrack},
		{"ontrack", linear.ProjectUpdateHealthTypeOntrack},
		{"at-risk", linear.ProjectUpdateHealthTypeAtrisk},
		{"atRisk", linear.ProjectUpdateHealthTypeAtrisk},
		{"off-track", linear.ProjectUpdateHealthTypeOfftrack},
		{"OFF-TRACK", linear.ProjectUpdateHealthTypeOfftrack},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Validate(tt.input)
			if err != nil {
				t.Fatalf("Validate(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("Validate(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidate_Invalid(t *testing.T) {
	for _, input := range []string{"green", "yellow", "", "on track"} {
		t.Run(input, func(t *testing.T) {
			if _, err := Validate(input); err == nil {
				t.Fatalf("Validate(%q) expected error", input)
			} else if !strings.Contains(err.Error(), Values) {
				t.Errorf("error should list valid values, got: %v", err)
			}
		})
	}
}
