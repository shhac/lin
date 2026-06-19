package mappers

// setIfNotEmpty assigns value to m[key] only when value is non-empty, so output
// maps omit blank optional fields rather than carrying empty strings.
func setIfNotEmpty(m map[string]any, key, value string) {
	if value != "" {
		m[key] = value
	}
}

// setIfNotNil assigns the dereferenced *v to m[key] only when v is non-nil, so
// output maps omit absent optional fields.
func setIfNotNil[T any](m map[string]any, key string, v *T) {
	if v != nil {
		m[key] = *v
	}
}
