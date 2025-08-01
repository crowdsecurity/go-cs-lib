package csyaml

import (
	"errors"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
)

// GetDocumentKeys reads all YAML documents from r and for each one
// returns a slice of its top-level keys, in order.
//
// Non-mapping documents yield an empty slice. Duplicate keys
// are not allowed and return an error.
func GetDocumentKeys(r io.Reader) ([][]string, error) {
	// Decode into Go types, but force mappings into MapSlice
	dec := yaml.NewDecoder(r, yaml.UseOrderedMap())

	allKeys := make([][]string, 0)

	idx := -1

	for {
		var raw any

		idx++

		if err := dec.Decode(&raw); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("position %d: %s", idx, yaml.FormatError(err, false, false))
		}
		keys := []string{}

		// Only mapping nodes become MapSlice with UseOrderedMap()
		if ms, ok := raw.(yaml.MapSlice); ok {
			for _, item := range ms {
				// Key is interface{}—here we expect strings
				if ks, ok := item.Key.(string); ok {
					keys = append(keys, ks)
				} else {
					// fallback to string form of whatever it is
					keys = append(keys, fmt.Sprint(item.Key))
				}
			}
		}

		allKeys = append(allKeys, keys)
	}

	return allKeys, nil
}
