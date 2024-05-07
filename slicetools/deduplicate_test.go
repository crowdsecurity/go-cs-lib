package slicetools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUniqueIntegers tests the UniqueCopy and Deduplicate functions with integer slices.
// Also test whether the original slice is modified or not.
func TestUniqueIntegers(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  []int
	}{
		{"Nil slice", nil, nil},
		{"Empty slice", []int{}, []int{}},
		// In test values, the subslice must differ from the original slice's backing array
		{"Unique integers", []int{1, 2, 1, 3, 2, 1}, []int{1, 2, 3}},
		{"Repeated integers", []int{4, 5, 5, 6, 4, 6, 5}, []int{4, 5, 6}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := UniqueCopy(tc.slice)
			assert.Equal(t, tc.want, got)

			if len(got) > 0 {
				assert.NotEqual(t, got, tc.slice[:len(got)], "UniqueCopy should not modify the original slice!")
			}

			got = Deduplicate(tc.slice)
			assert.Equal(t, tc.want, got)

			if len(got) > 0 {
				assert.Equal(t, got, tc.slice[:len(got)], "Deduplicate should modify the original slice!")
			}

			for i := len(got); i < len(tc.slice); i++ {
				assert.Equal(t, 0, tc.slice[i], "Deduplicate should zero out the remaining elements!")
			}
		})
	}
}

func TestUniqueStrings(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		want  []string
	}{
		{"Empty slice", []string{}, []string{}},
		{"Unique strings", []string{"a", "b", "a", "c", "b"}, []string{"a", "b", "c"}},
		{"Repeated strings", []string{"hello", "world", "hello"}, []string{"hello", "world"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := UniqueCopy(tc.slice)
			assert.Equal(t, tc.want, got)
			got = Deduplicate(tc.slice)
			assert.Equal(t, tc.want, got)
		})
	}
}

type Point struct {
	X int
	Y int
}

func TestUniquePoints(t *testing.T) {
	points := []Point{
		{X: 1, Y: 2},
		{X: 3, Y: 4},
		{X: 1, Y: 2}, // Duplicate of the first point
		{X: 5, Y: 6},
	}

	want := []Point{
		{X: 1, Y: 2},
		{X: 3, Y: 4},
		{X: 5, Y: 6},
	}

	got := UniqueCopy(points)
	assert.Equal(t, want, got)
	got = Deduplicate(points)
	assert.Equal(t, want, got)
}
