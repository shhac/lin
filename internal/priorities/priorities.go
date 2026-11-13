package priorities

import (
	"strings"

	apierrors "github.com/shhac/lin/internal/errors"
)

var Map = map[string]int{
	"none":   0,
	"urgent": 1,
	"high":   2,
	"medium": 3,
	"low":    4,
}

const Values = "none | urgent | high | medium | low"

// Resolve returns the priority number for a name, or -1 if unrecognized.
func Resolve(input string) (int, bool) {
	p, ok := Map[strings.ToLower(input)]
	return p, ok
}

// Validate resolves a priority name, or returns the invalid-choice error every
// command reports for one.
func Validate(input string) (int, error) {
	p, ok := Resolve(input)
	if !ok {
		return 0, apierrors.Newf(apierrors.FixableByAgent, "invalid priority: %q, valid values: %s", input, Values)
	}
	return p, nil
}
