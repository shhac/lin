// Package projecthealth validates and normalizes project-update health values.
package projecthealth

import (
	"fmt"
	"strings"

	"github.com/shhac/lin/internal/linear"
)

// Values is the human-facing list of accepted health inputs.
const Values = "on-track | at-risk | off-track"

// byInput maps accepted spellings (hyphenated, camelCase, and bare) to the
// Linear ProjectUpdateHealthType enum.
var byInput = map[string]linear.ProjectUpdateHealthType{
	"on-track":  linear.ProjectUpdateHealthTypeOntrack,
	"ontrack":   linear.ProjectUpdateHealthTypeOntrack,
	"at-risk":   linear.ProjectUpdateHealthTypeAtrisk,
	"atrisk":    linear.ProjectUpdateHealthTypeAtrisk,
	"off-track": linear.ProjectUpdateHealthTypeOfftrack,
	"offtrack":  linear.ProjectUpdateHealthTypeOfftrack,
}

// Validate resolves a human-supplied health value (case-insensitive,
// hyphen-optional) to the Linear enum, or returns an error listing valid values.
func Validate(input string) (linear.ProjectUpdateHealthType, error) {
	if h, ok := byInput[strings.ToLower(input)]; ok {
		return h, nil
	}
	return "", fmt.Errorf("unknown project health: %q, valid values: %s", input, Values)
}
