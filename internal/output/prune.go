package output

import "strings"

func pruneEmpty(v any) any {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, v := range val {
			pruned := pruneEmpty(v)
			if pruned == nil {
				continue
			}
			// pruneEmpty preserves empty slices at top level; inside maps,
			// drop them too so pruned objects don't keep stub keys.
			if arr, ok := pruned.([]any); ok && len(arr) == 0 {
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
			pruned := pruneEmpty(v)
			if pruned != nil {
				out = append(out, pruned)
			}
		}
		return out
	case string:
		if strings.TrimSpace(val) == "" {
			return nil
		}
		return val
	default:
		return v
	}
}
