package slicetools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestBackwardInts(t *testing.T) {
	ints := []int{10, 20, 30, 40, 50}
	expected := []int{50, 40, 30, 20, 10}
	expectedIdx := []int{4, 3, 2, 1, 0}

	result := []int{}
	resultIdx := []int{}

	for i, value := range Backward(ints) {
		result = append(result, value)
		resultIdx = append(resultIdx, i)
	}

	assert.Equal(t, expected, result)
	assert.Equal(t, expectedIdx, resultIdx)
}

// TestBackwardStrings tests the Backward function with a slice of strings.
func TestBackwardStrings(t *testing.T) {
	strs := []string{"apple", "banana", "cherry", "date"}
	expected := []string{"date", "cherry", "banana", "apple"}
	expectedIdx := []int{3, 2, 1, 0}

	result := []string{}
	resultIdx := []int{}

	for i, value := range Backward(strs) {
		result = append(result, value)
		resultIdx = append(resultIdx, i)
	}

	assert.Equal(t, expected, result)
	assert.Equal(t, expectedIdx, resultIdx)
}

func TestBackwardEarlyTermination(t *testing.T) {
	ints := []int{10, 20, 30, 40, 50}
	expected := []int{50, 40, 30} // Stop after reaching 30

	var result []int

	for _, value := range Backward(ints) {
		result = append(result, value)
		if value == 30 {
			break
		}
	}

	assert.Equal(t, expected, result)
}

func TestBackwardEmptySlice(t *testing.T) {
	ints := []int{}

	result := []int{}
	resultIdx := []int{}

	for i, value := range Backward(ints) {
		result = append(result, value)
		resultIdx = append(resultIdx, i)
	}

	assert.Empty(t, result)
	assert.Empty(t, resultIdx)
}
