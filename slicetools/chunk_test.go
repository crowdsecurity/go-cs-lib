package slicetools_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/crowdsecurity/go-cs-lib/slicetools"
)

func TestChunks(t *testing.T) {
	testCases := []struct {
		name      string
		items     []int
		chunkSize int
		expected  [][]int
	}{
		{ "empty slice, chunk size 2", []int{}, 2, [][]int{}},
		{"1 element, chunk size 2",    []int{1}, 2, [][]int{{1}}},
		{"empty slice, chunk size 0",  []int{}, 0, [][]int{}},
		{"5 elements, chunk size 2",   []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{"5 elements, chunk size 3",   []int{1, 2, 3, 4, 5}, 3, [][]int{{1, 2, 3}, {4, 5}}},
		{"5 elements, chunk size 4",   []int{1, 2, 3, 4, 5}, 5, [][]int{{1, 2, 3, 4, 5}}},
		{"5 elements, chunk size 6",   []int{1, 2, 3, 4, 5}, 6, [][]int{{1, 2, 3, 4, 5}}},
		{"chunk size 0 = don't chunk", []int{1, 2, 3, 4, 5}, 0, [][]int{{1, 2, 3, 4, 5}}},
		{"look ma, no sorting",        []int{1, 2, 4, 1, 5}, 2, [][]int{{1, 2}, {4, 1}, {5}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := slicetools.Chunks(tc.items, tc.chunkSize)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
