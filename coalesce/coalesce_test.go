package coalesce

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crowdsecurity/go-cs-lib/ptr"
)

func TestString(t *testing.T) {
	tests := []struct{have []string; want string}{
		{[]string{}, ""},
		{[]string{"", "", ""}, ""},
		{[]string{"", "", "c"}, "c"},
		{[]string{"", "b", ""}, "b"},
		{[]string{"a", "", ""}, "a"},
		{[]string{"a", "b", "c"}, "a"},
	}

	for _, tc := range tests {
		got := String(tc.have...)
		assert.Equal(t, tc.want, got)
	}
}

func TestInt(t *testing.T) {
	tests := []struct{have []int; want int}{
		{[]int{}, 0},
		{[]int{0, 0, 0}, 0},
		{[]int{0, 0, 3}, 3},
		{[]int{0, 2, 0}, 2},
		{[]int{1, 0, 0}, 1},
		{[]int{1, 2, 3}, 1},
	}

	for _, tc := range tests {
		got := Int(tc.have...)
		assert.Equal(t, tc.want, got)
	}
}

func TestNotNil(t *testing.T) {
	tests := []struct{have []*int; want *int}{
		{[]*int{}, nil},
		{[]*int{nil, nil, nil}, nil},
		{[]*int{nil, nil, ptr.Of(3)}, ptr.Of(3)},
		{[]*int{nil, ptr.Of(0), ptr.Of(3)}, ptr.Of(0)},
		{[]*int{nil, ptr.Of(2), nil}, ptr.Of(2)},
		{[]*int{ptr.Of(1), nil, nil}, ptr.Of(1)},
		{[]*int{ptr.Of(1), ptr.Of(2), ptr.Of(3)}, ptr.Of(1)},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.have), func(t *testing.T) {
			got := NotNil(tc.have...)
			if tc.want == nil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, *tc.want, *got)
		})
	}
}

/*
func TestNotEmptyOrNil(t *testing.T) {
	tests := []struct{have []*int; want *int}{
		{[]*int{}, nil},
		{[]*int{nil, nil, nil}, nil},
		{[]*int{nil, nil, ptr.Of(3)}, ptr.Of(3)},
		{[]*int{nil, ptr.Of(0), ptr.Of(3)}, ptr.Of(3)},
		{[]*int{nil, ptr.Of(2), nil}, ptr.Of(2)},
		{[]*int{ptr.Of(1), nil, nil}, ptr.Of(1)},
		{[]*int{ptr.Of(1), ptr.Of(2), ptr.Of(3)}, ptr.Of(1)},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%v", tc.have), func(t *testing.T) {
			got := NotEmptyOrNil(tc.have...)
			if tc.want == nil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, *tc.want, *got)
		})
	}
}
*/
