package csyaml_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crowdsecurity/go-cs-lib/cstest"
	"github.com/crowdsecurity/go-cs-lib/csyaml"
)

func TestMergeYAML(t *testing.T) {
	tests := []struct {
		name    string
		inputs  []string
		want    string
		wantErr string
	}{
		{
			name:   "single doc passes through",
			inputs: []string{"a: 1\nb: 2\n"},
			want:   "a: 1\nb: 2\n",
		},
		{
			name: "merge maps deep",
			inputs: []string{
				"one: 1\ntwo: 2\n",
				"two: 20\nthree: 3\n",
			},
			want: "one: 1\ntwo: 20\nthree: 3\n",
		},
		{
			name: "sequence replaced",
			inputs: []string{
				"list: [1,2,3]\n",
				"list: [4,5]\n",
			},
			want: "list:\n- 4\n- 5\n",
		},
		{
			name: "scalar override",
			inputs: []string{
				"foo: bar\n",
				"foo: baz\n",
			},
			want: "foo: baz\n",
		},
		{
			name:    "type mismatch error",
			inputs:  []string{"foo: 1\n", "foo:\n  - a\n"},
			wantErr: "cannot merge sequence into scalar",
		},
		{
			name:    "invalid yaml error",
			inputs:  []string{"list: [1,2,"},
			wantErr: "decoding document 0: [1:7] sequence end token ']' not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// build byte slices
			var bs [][]byte
			for _, s := range tc.inputs {
				bs = append(bs, []byte(s))
			}

			buf, err := csyaml.Merge(bs)

			cstest.RequireErrorContains(t, err, tc.wantErr)
			if tc.wantErr != "" {
				require.Nil(t, buf)
			} else {
				assert.Equal(t, tc.want, buf.String())
			}
		})
	}
}
