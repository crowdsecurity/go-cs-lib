package maptools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedKeys(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want []string
	}{
		{
			name: "empty map",
			m:    map[string]int{},
			want: []string{},
		},
		{
			name: "single element",
			m:    map[string]int{"a": 1},
			want: []string{"a"},
		},
		{
			name: "multiple elements",
			m:    map[string]int{"b": 2, "a": 1, "c": 3},
			want: []string{"a", "b", "c"},
		},
		{
			name: "elements with same values",
			m:    map[string]int{"b": 2, "a": 2, "c": 2},
			want: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, SortedKeys(tt.m))
		})
	}
}

func TestSortedKeysNoCase(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want []string
	}{
		{
			name: "empty map",
			m:    map[string]int{},
			want: []string{},
		},
		{
			name: "single element",
			m:    map[string]int{"a": 1},
			want: []string{"a"},
		},
		{
			name: "multiple elements with different cases",
			m:    map[string]int{"b": 2, "A": 1, "C": 3},
			want: []string{"A", "b", "C"},
		},
		{
			name: "elements with same values and different cases",
			m:    map[string]int{"b": 2, "A": 2, "c": 2},
			want: []string{"A", "b", "c"},
		},
		{
			name: "mixed case elements",
			m:    map[string]int{"Banana": 1, "apple": 2, "Cherry": 3, "banana": 4, "Apple": 5, "cherry": 6},
			want: []string{"Apple", "apple", "Banana", "banana", "Cherry", "cherry"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, SortedKeysNoCase(tt.m))
		})
	}
}
