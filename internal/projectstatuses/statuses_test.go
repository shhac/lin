package projectstatuses

import (
	"strings"
	"testing"
)

func TestValidate_ValidStatuses(t *testing.T) {
	for _, status := range List {
		t.Run(status, func(t *testing.T) {
			got, err := Validate(status)
			if err != nil {
				t.Fatalf("Validate(%q) unexpected error: %v", status, err)
			}
			if got != status {
				t.Errorf("Validate(%q) = %q, want %q", status, got, status)
			}
		})
	}
}

func TestValidate_CaseInsensitive(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Backlog", "backlog"},
		{"PLANNED", "planned"},
		{"Started", "started"},
		{"PAUSED", "paused"},
		{"Completed", "completed"},
		{"CANCELED", "canceled"},
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
	tests := []string{"active", "done", "archived", "", "in_progress"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := Validate(input)
			if err == nil {
				t.Fatalf("Validate(%q) expected error", input)
			}
			if !strings.Contains(err.Error(), "unknown project status") {
				t.Errorf("unexpected error message: %v", err)
			}
			if !strings.Contains(err.Error(), Values) {
				t.Errorf("error should include valid values, got: %v", err)
			}
		})
	}
}
