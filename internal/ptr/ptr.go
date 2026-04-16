package ptr

// To returns a pointer to the given value.
func To[T any](v T) *T { return &v }

// TrueOrNil returns &true if b is true, nil otherwise.
// Useful for optional bool flags where false means "omit", not "send false".
func TrueOrNil(b bool) *bool {
	if !b {
		return nil
	}
	return To(true)
}

// Deref returns the value pointed to, or the zero value if nil.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
