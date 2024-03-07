package slicetools

// UniqueCopy filters a slice of any comparable type, returning a slice of unique elements.
// Keeps the order of the original slice.
func UniqueCopy[T comparable](slice []T) []T {
	if slice == nil {
		return nil
	}
	seen := make(map[T]struct{})
	ret := make([]T, 0, len(slice))
	for _, value := range slice {
		if _, ok := seen[value]; !ok {
			seen[value] = struct{}{}
			ret = append(ret, value)
		}
	}
	return ret
}

// Deduplicate filters a slice of any comparable type, removing duplicate elements in place.
// Keeps the order of the original slice.
// The original slice can't be used after calling this function because it's likely to have a different length.
func Deduplicate[T comparable](slice []T) []T {
	if slice == nil {
		return nil
	}
	seen := make(map[T]struct{})
	j := 0
	for _, value := range slice {
		if _, ok := seen[value]; !ok {
			seen[value] = struct{}{}
			slice[j] = value
			j++
		}
	}

	// Zero elements from j to len(slice), like stdlib methods
	var zero T // Default zero value for type T
	for i := j; i < len(slice); i++ {
		slice[i] = zero
	}

	return slice[:j]
}
