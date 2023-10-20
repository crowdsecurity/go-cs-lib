package slicetools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcat(t *testing.T) {
	tests := []struct {
		name   string
		slices [][]int
		want   []int
	}{
		{
			name:   "all empty slices",
			slices: [][]int{{}, {}, {}},
			want:   []int{},
		},
		{
			name:   "mixed empty and non-empty slices",
			slices: [][]int{{}, {1, 2}, {}, {3, 4, 5}, {}},
			want:   []int{1, 2, 3, 4, 5},
		},
		{
			name:   "non-empty slices",
			slices: [][]int{{6, 7}, {8, 9, 10}},
			want:   []int{6, 7, 8, 9, 10},
		},
		{
			name:   "single empty slice",
			slices: [][]int{{}},
			want:   []int{},
		},
		{
			name:   "no slices",
			slices: [][]int{},
			want:   []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Concat(tc.slices...)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestConcatStrings(t *testing.T) {
	tests := []struct {
		name   string
		slices [][]string
		want   []string
	}{
		{
			name:   "all empty slices of strings",
			slices: [][]string{{}, {}, {}},
			want:   []string{},
		},
		{
			name:   "mixed empty and non-empty slices of strings",
			slices: [][]string{{}, {"a", "b"}, {}, {"c", "d", "e"}, {}},
			want:   []string{"a", "b", "c", "d", "e"},
		},
		{
			name:   "non-empty slices of strings",
			slices: [][]string{{"f", "g"}, {"h", "i", "j"}},
			want:   []string{"f", "g", "h", "i", "j"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Concat(tc.slices...)
			assert.Equal(t, tc.want, got)
		})
	}
}
