package slicetools

// Concat concatenates multiple slices and returns the result.
func Concat[T any](slices...[]T) []T {
	tot := 0
	for _, s := range slices {
		tot += len(s)
	}

	ret := make([]T, tot)

	i := 0
	for _, s := range slices {
		i += copy(ret[i:], s)
	}

	return ret
}
