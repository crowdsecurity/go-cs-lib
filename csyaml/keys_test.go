package csyaml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crowdsecurity/go-cs-lib/cstest"
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
			got, err := GetDocumentKeys(r)
			cstest.RequireErrorContains(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
