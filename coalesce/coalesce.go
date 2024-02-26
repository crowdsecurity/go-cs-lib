package coalesce

// String returns the first non-empty string, or ""
func String(s ...string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}

	return ""
}

// Int returns the first non-zero value, or zero
func Int(s ...int) int {
	for _, v := range s {
		if v != 0 {
			return v
		}
	}

	return 0
}

// NotNil returns the first non-nil pointer, or nil
func NotNil[T any](args ...*T) *T {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}

	return nil
}

/* Is this useful? Here it is, just in case

// NotEmpty returns the first non-nil pointer to a non-zero value, or nil
func NotEmptyOrNil[T comparable](args ...*T) *T {
	var zero T
	for _, arg := range args {
		if arg != nil && *arg != zero {
			return arg
		}
	}

	return nil
}
*/
