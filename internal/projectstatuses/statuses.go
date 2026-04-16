package projectstatuses

import (
	"fmt"
	"strings"
)

var List = []string{"backlog", "planned", "started", "paused", "completed", "canceled"}

const Values = "backlog | planned | started | paused | completed | canceled"

// Validate checks if the input matches a known project status (case-insensitive).
// Returns the normalized (lowercase) status or an error.
func Validate(input string) (string, error) {
	lower := strings.ToLower(input)
	for _, s := range List {
		if strings.EqualFold(s, lower) {
			return s, nil
		}
	}
	return "", fmt.Errorf("unknown project status: %q, valid values: %s", input, Values)
}
