package testutil

import "encoding/json"

// MustJSON marshals v to pretty-printed JSON, panicking on error.
func MustJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic("MustJSON: " + err.Error())
	}
	return string(b)
}
