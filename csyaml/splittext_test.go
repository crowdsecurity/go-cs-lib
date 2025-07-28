package csyaml_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crowdsecurity/go-cs-lib/cstest"
	"github.com/crowdsecurity/go-cs-lib/csyaml"
)

func TestSplitDocumentsText(t *testing.T) {
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
			want:  [][]byte{[]byte("\"scalar\"\n")},
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
				[]byte("---\na: 1\nb: 2\n"),
				[]byte("---\n- 1\n- 2\n"),
				[]byte("---\n\"scalar\"\n"),
			},
		},
		{
			name:  "empty input",
			input: "",
			want:  [][]byte(nil),
		},
		{
			name:  "invalid YAML",
			input: "list: [1, 2,",
			want:  [][]byte{[]byte("list: [1, 2,\n")},
		},
		{
			name: "preserve comments",
			input: `# comment 1
a: 1
# comment 2
b: 2
---
# comment 3
- 1
# comment 4
- 2
# comment 5
---
# comment 6
"scalar"
# comment 7
`,
			want: [][]byte{
				[]byte("# comment 1\na: 1\n# comment 2\nb: 2\n"),
				[]byte("---\n# comment 3\n- 1\n# comment 4\n- 2\n# comment 5\n"),
				[]byte("---\n# comment 6\n\"scalar\"\n# comment 7\n"),
			},
		},
		{
			name: "tricky separator",
			input: `---
text: |
  This is a multi-line string.
  It includes a line that looks like a document separator:
  ---
  But it's just part of the string.
---
key: value
`,
			want: [][]byte{
				[]byte("---\ntext: |\n  This is a multi-line string.\n  It includes a line that looks like a document separator:\n"),
				[]byte("  ---\n  But it's just part of the string.\n"),
				[]byte("---\nkey: value\n"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			docs, err := csyaml.SplitDocumentsText(r)
			cstest.RequireErrorContains(t, err, tc.wantErr)
			assert.Equal(t, tc.want, docs)
		})
	}
}
