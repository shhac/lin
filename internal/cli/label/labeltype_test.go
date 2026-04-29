package label

import (
	"strings"
	"testing"
)

func TestValidateType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"issue is valid", "issue", false},
		{"project is valid", "project", false},
		{"empty string is invalid", "", true},
		{"unknown is invalid", "task", true},
		{"uppercase is invalid (case-sensitive)", "ISSUE", true},
		{"whitespace is invalid", " issue", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateType(tt.input)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error for %q, got nil", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil for %q, got %v", tt.input, err)
			}
			if tt.wantErr && !strings.Contains(err.Error(), "--type must be") {
				t.Errorf("error message should mention `--type must be`, got: %v", err)
			}
		})
	}
}

func TestRejectTeamForProject(t *testing.T) {
	tests := []struct {
		name     string
		typeFlag string
		teamFlag string
		wantErr  bool
		wantHint bool
	}{
		{"project + team is rejected", "project", "ENG", true, true},
		{"project + empty team is allowed", "project", "", false, false},
		{"issue + team is allowed", "issue", "ENG", false, false},
		{"issue + empty team is allowed", "issue", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rejectTeamForProject(tt.typeFlag, tt.teamFlag)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if tt.wantErr && !strings.Contains(err.Error(), "--team is not valid with --type=project") {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateLabelTypeFlags_StopsAtFirstError(t *testing.T) {
	// Bad --type and team-on-project are both errors; the type check fires first.
	err := validateLabelTypeFlags("bogus", "ENG")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "--type must be") {
		t.Errorf("expected type error first, got: %v", err)
	}
}

func TestValidateLabelTypeFlags_AllValid(t *testing.T) {
	if err := validateLabelTypeFlags("issue", "ENG"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := validateLabelTypeFlags("project", ""); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
