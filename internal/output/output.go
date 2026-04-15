package output

import (
	"encoding/json"
	"fmt"
	"os"
)

type Pagination struct {
	HasMore    bool   `json:"hasMore"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type PaginatedResult struct {
	Items      any         `json:"items"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func PrintJSON(data any) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	var decoded any
	if err := json.Unmarshal(b, &decoded); err != nil {
		return
	}
	decoded = pruneEmpty(decoded)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	_ = enc.Encode(decoded)
}

func PrintPaginated(items any, pageInfo *Pagination) {
	result := PaginatedResult{Items: items}
	if pageInfo != nil && pageInfo.HasMore {
		result.Pagination = pageInfo
	}
	PrintJSON(result)
}

func PrintError(msg string) {
	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(map[string]string{"error": msg})
	os.Exit(1)
}

func PrintErrorf(format string, args ...any) {
	PrintError(fmt.Sprintf(format, args...))
}

func pruneEmpty(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, v := range val {
			if v == nil {
				continue
			}
			pruned := pruneEmpty(v)
			if isEmpty(pruned) {
				continue
			}
			out[k] = pruned
		}
		if len(out) == 0 {
			return nil
		}
		return out
	case []any:
		out := make([]any, 0, len(val))
		for _, v := range val {
			out = append(out, pruneEmpty(v))
		}
		return out
	case string:
		if val == "" {
			return nil
		}
		return val
	default:
		return v
	}
}

func isEmpty(v any) bool {
	switch val := v.(type) {
	case nil:
		return true
	case string:
		return val == ""
	case map[string]any:
		return len(val) == 0
	case []any:
		return len(val) == 0
	default:
		return false
	}
}
