package projectstatuses

import (
	"strings"

	apierrors "github.com/shhac/lin/internal/errors"
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
	return "", apierrors.Newf(apierrors.FixableByAgent, "unknown project status: %q, valid values: %s", input, Values)
}
