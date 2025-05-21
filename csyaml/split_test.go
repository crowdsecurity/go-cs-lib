package csyaml_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crowdsecurity/go-cs-lib/cstest"
	"github.com/crowdsecurity/go-cs-lib/csyaml"
)

func TestSplitDocuments(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    [][]byte
		wantErr string
	}{
		{
			name:  "single mapping",
			input: "a: 1\nb: 2\n",
			want:  [][]byte{[]byte("a: 1\nb: 2\n")},
		},
		{
			name:  "sequence doc",
			input: "- 1\n- 2\n",
			want:  [][]byte{[]byte("- 1\n- 2\n")},
		},
		{
			name:  "scalar doc",
			input: "\"scalar\"\n",
			want:  [][]byte{[]byte("scalar\n")},
		},
		{
			name: "multiple documents",
			input: `---
a: 1
b: 2
---
- 1
- 2
---
"scalar"
`,
			want: [][]byte{
				[]byte("a: 1\nb: 2\n"),
				[]byte("- 1\n- 2\n"),
				[]byte("scalar\n"),
			},
		},
		{
			name:  "empty input",
			input: "",
			want:  [][]byte{},
		},
		{
			name:    "invalid YAML",
			input:   "list: [1, 2,",
			want:    nil,
			wantErr: "decoding document 0: [1:7] sequence end token ']' not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			docs, err := csyaml.SplitDocuments(r)
			cstest.RequireErrorContains(t, err, tc.wantErr)
			assert.Equal(t, tc.want, docs)
		})
	}
}
