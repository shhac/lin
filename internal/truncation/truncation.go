package truncation

import (
	"fmt"
	"strings"
)

const defaultMaxLength = 200
const ellipsis = "\u2026"

var (
	expandedFields map[string]bool
	expandAll      bool
	maxLength      = defaultMaxLength
)

var truncatableFields = map[string]bool{
	"description": true,
	"body":        true,
	"content":     true,
}

// Configure sets truncation state for the CLI invocation.
func Configure(opts ConfigOpts) {
	expandedFields = make(map[string]bool)
	expandAll = false
	maxLength = defaultMaxLength

	if opts.Full {
		expandAll = true
	}
	if opts.Expand != "" {
		for _, f := range strings.Split(opts.Expand, ",") {
			expandedFields[strings.ToLower(strings.TrimSpace(f))] = true
		}
	}
	if opts.MaxLength > 0 {
		maxLength = opts.MaxLength
	}
}

type ConfigOpts struct {
	Expand    string
	Full      bool
	MaxLength int
}

func shouldExpand(field string) bool {
	return expandAll || expandedFields[strings.ToLower(field)]
}

// Apply recursively truncates truncatable fields in the data structure.
// Works on the map[string]any representation produced by JSON marshal/unmarshal.
func Apply(data any) any {
	if data == nil {
		return nil
	}
	switch v := data.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for k, val := range v {
			if truncatableFields[k] {
				if s, ok := val.(string); ok {
					out[fmt.Sprintf("%sLength", k)] = len(s)
					if shouldExpand(k) || len(s) <= maxLength {
						out[k] = s
					} else {
						out[k] = s[:maxLength] + ellipsis
					}
					continue
				}
			}
			out[k] = Apply(val)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			out[i] = Apply(item)
		}
		return out
	default:
		return data
	}
}
