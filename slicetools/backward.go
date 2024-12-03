package slicetools

// Backward iterates over a slice in reverse order.
func Backward[E any](s []E) func(func(int, E) bool) {
	return func(yield func(int, E) bool) {
		for i := len(s)-1; i >= 0; i-- {
			if !yield(i, s[i]) {
				return
			}
		}
	}
}
