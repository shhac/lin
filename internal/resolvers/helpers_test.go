package resolvers

import (
	"reflect"
	"testing"
)

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty string is empty", "", nil},
		{"only spaces is empty", "   ", nil},
		{"only commas is empty", ",,,", nil},
		{"comma + spaces is empty", " , , ", nil},
		{"single value", "Bug", []string{"Bug"}},
		{"trims whitespace", "  Bug  ", []string{"Bug"}},
		{"two comma-separated", "Bug,Feature", []string{"Bug", "Feature"}},
		{"trims around commas", " Bug , Feature ", []string{"Bug", "Feature"}},
		{"leading comma is dropped", ",Bug", []string{"Bug"}},
		{"trailing comma is dropped", "Bug,", []string{"Bug"}},
		{"empty middle token is dropped", "Bug,,Feature", []string{"Bug", "Feature"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitAndTrim(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitAndTrim(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
