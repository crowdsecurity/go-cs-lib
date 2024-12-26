package maptools

import (
	"sort"
	"strings"
)

func SortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func SortedKeysNoCase[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		li, lj := strings.ToLower(keys[i]), strings.ToLower(keys[j])
		if li == lj {
			// differ only by case, sort by original key
			// will have uppercase first
			return keys[i] < keys[j]
		}

		return li < lj
	})

	return keys
}
