package csyaml_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crowdsecurity/go-cs-lib/cstest"
	"github.com/crowdsecurity/go-cs-lib/csyaml"
)

func TestCollectTopLevelKeys(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    [][]string
		wantErr string
	}{
		{
			name:  "single mapping",
			input: "a: 1\nb: 2\n",
			want:  [][]string{{"a", "b"}},
		},
		{
			name:    "duplicate keys mapping",
			input:   "a: 1\nb: 2\na: 3\n",
			want:    nil,
			wantErr: `position 0: [3:1] mapping key "a" already defined at [1:1]`,
		},
		{
			name: "multiple documents",
			input: `---
a: 1
b: 2
---
- 1
---
c: 1
b: 2
---
"scalar"
`,
			want: [][]string{{"a", "b"}, {}, {"c", "b"}, {}},
		},
		{
			name:  "empty input",
			input: "",
			want:  [][]string{},
		},
		{
			name:    "invalid YAML",
			input:   "list: [1, 2,",
			want:    nil,
			wantErr: "position 0: [1:7] sequence end token ']' not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			got, err := csyaml.GetDocumentKeys(r)
			cstest.RequireErrorContains(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
