package priorities

import (
	"strings"
	"testing"
)

func TestResolve_ValidPriorities(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"none", 0},
		{"urgent", 1},
		{"high", 2},
		{"medium", 3},
		{"low", 4},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, ok := Resolve(tt.input)
			if !ok {
				t.Fatalf("Resolve(%q) returned ok=false", tt.input)
			}
			if got != tt.want {
				t.Errorf("Resolve(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestResolve_CaseInsensitive(t *testing.T) {
	tests := []string{"High", "HIGH", "hIgH", "URGENT", "Low"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, ok := Resolve(input)
			if !ok {
				t.Errorf("Resolve(%q) should succeed (case-insensitive)", input)
			}
		})
	}
}

func TestResolve_Invalid(t *testing.T) {
	tests := []string{"critical", "blocker", "", "p1", "highest"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, ok := Resolve(input)
			if ok {
				t.Errorf("Resolve(%q) should return ok=false", input)
			}
		})
	}
}

func TestMap_AllPresent(t *testing.T) {
	expected := []string{"none", "urgent", "high", "medium", "low"}
	for _, name := range expected {
		if _, ok := Map[name]; !ok {
			t.Errorf("Map missing key %q", name)
		}
	}
	if len(Map) != len(expected) {
		t.Errorf("Map has %d entries, expected %d", len(Map), len(expected))
	}
}

func TestValidate(t *testing.T) {
	if p, err := Validate("High"); err != nil || p != 2 {
		t.Errorf("Validate(High) = %d, %v", p, err)
	}
	_, err := Validate("bogus")
	if err == nil || !strings.Contains(err.Error(), `invalid priority: "bogus", valid values: none | urgent | high | medium | low`) {
		t.Errorf("Validate(bogus) err = %v", err)
	}
}
