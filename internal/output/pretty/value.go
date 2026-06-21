package pretty

// Helpers for reading the raw map[string]any a mapper produces. Values keep
// their Go types (mappers run before any JSON round-trip), so strings may be
// string or *string, numbers float64/int or their pointers.

// Str reads a string field, dereferencing *string and returning "" for nil.
func Str(m map[string]any, key string) string {
	switch v := m[key].(type) {
	case string:
		return v
	case *string:
		if v != nil {
			return *v
		}
	}
	return ""
}

// Num reads a numeric field as float64, dereferencing pointers. ok is false when
// the field is absent or nil.
func Num(m map[string]any, key string) (float64, bool) {
	switch v := m[key].(type) {
	case float64:
		return v, true
	case *float64:
		if v != nil {
			return *v, true
		}
	case int:
		return float64(v), true
	case *int:
		if v != nil {
			return float64(*v), true
		}
	}
	return 0, false
}

// Int reads a numeric field as int (truncating), 0 when absent.
func Int(m map[string]any, key string) int {
	if f, ok := Num(m, key); ok {
		return int(f)
	}
	return 0
}

// Submap reads a nested object field, nil when absent.
func Submap(m map[string]any, key string) map[string]any {
	if v, ok := m[key].(map[string]any); ok {
		return v
	}
	return nil
}

// MapSlice reads an array-of-objects field. It accepts both []map[string]any
// (the mapper's native shape) and []any (post-JSON), nil when absent.
func MapSlice(m map[string]any, key string) []map[string]any {
	switch v := m[key].(type) {
	case []map[string]any:
		return v
	case []any:
		out := make([]map[string]any, 0, len(v))
		for _, e := range v {
			if mm, ok := e.(map[string]any); ok {
				out = append(out, mm)
			}
		}
		return out
	}
	return nil
}
