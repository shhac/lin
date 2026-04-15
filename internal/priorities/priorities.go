package priorities

import "strings"

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
