package csyaml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crowdsecurity/go-cs-lib/cstest" // adjust this import to your package
)

func TestIsEmptyYAML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr string
	}{
		{
			name:  "empty document",
			input: ``,
			want:  true,
		},
		{
			name:  "just a key",
			input: "foo:",
			want:  false,
		},
		{
			name:  "just newline",
			input: "\n",
			want:  true,
		},
		{
			name:  "just comment",
			input: "# only a comment",
			want:  true,
		},
		{
			name:  "comments and empty lines",
			input: "# only a comment\n\n# another one\n\n",
			want:  true,
		},
		{
			name:  "empty doc with separator",
			input: "---",
			want:  true,
		},
		{
			name:  "empty mapping",
			input: "{}",
			want:  false,
		},
		{
			name:  "empty sequence",
			input: "[]",
			want:  false,
		},
		{
			name:  "non-empty mapping",
			input: "foo: bar",
			want:  false,
		},
		{
			name:  "non-empty sequence",
			input: "- 1\n- 2",
			want:  false,
		},
		{
			name:  "non-empty scalar",
			input: "hello",
			want:  false,
		},
		{
			name:  "empty scalar",
			input: "''",
			want:  false,
		},
		{
			name:  "explicit nil",
			input: "null",
			want:  false,
		},
		{
			name:    "malformed YAML",
			input:   "foo: [1,",
			wantErr: "[1:6] sequence end token ']' not found",
		},
		{
			name:  "multiple empty documents",
			input: "---\n---\n---\n#comment",
			want:  true,
		},
		{
			name:  "second document is not empty",
			input: "---\nfoo: bar\n---\n#comment",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := IsEmptyYAML(strings.NewReader(tc.input))

			cstest.RequireErrorContains(t, err, tc.wantErr)

			if tc.wantErr != "" {
				return
			}

			assert.Equal(t, tc.want, got)
		})
	}
}
