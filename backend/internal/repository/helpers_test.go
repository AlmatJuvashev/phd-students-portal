package repository

// strPtr returns a pointer to the given string.
func strPtr(s string) *string {
	return &s
}

// toPtr returns a pointer to the given value.
func toPtr[T any](v T) *T {
	return &v
}
