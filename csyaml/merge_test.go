package csyaml_test

import (
	"strings"
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
			wantErr: "can't merge a sequence into a scalar",
		},
		{
			name:    "invalid yaml error",
			inputs:  []string{"ref: *foo\n"}, // undefined alias
			wantErr: `decoding document 0: [1:7] could not find alias "foo"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bs [][]byte
			for _, s := range tc.inputs {
				bs = append(bs, []byte(s))
			}

			buf, err := csyaml.Merge(bs)
			cstest.RequireErrorContains(t, err, tc.wantErr)

			if tc.wantErr != "" {
				require.Nil(t, buf)
				return
			}

			require.NotNil(t, buf)
			assert.Equal(t, tc.want, buf.String())
		})
	}
}

func TestEmptyVsNilSources(t *testing.T) {
	tests := []struct {
		desc    string
		sources [][]byte
		expect  string
	}{
		{"empty base", [][]byte{nil, []byte("foo: bar\n")}, "foo: bar\n"},
		{"empty override", [][]byte{[]byte("foo: bar\n"), nil}, "foo: bar\n"},
		{"both empty", [][]byte{nil, nil}, ""},
		{"null base", [][]byte{[]byte("~\n"), []byte("foo: bar\n")}, "foo: bar\n"},
		{"explicit null override", [][]byte{[]byte("foo: bar\n"), []byte("~\n")}, "null\n"},
		{"empty base & null override", [][]byte{nil, []byte("~\n")}, "null\n"},
		{"null base & empty override", [][]byte{[]byte("~\n"), nil}, "null\n"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			merged, err := csyaml.Merge(tt.sources)
			require.NoError(t, err)
			assert.Equal(t, tt.expect, merged.String())
		})
	}
}

func TestDuplicateKeyError(t *testing.T) {
	src := []byte("{foo: bar, foo: baz}")
	_, err := csyaml.Merge([][]byte{src})
	cstest.RequireErrorContains(t, err, `decoding document 0: [1:12] mapping key "foo" already defined at [1:2]`)
}

func TestTabsInSource(t *testing.T) {
	src := []byte("foo:\n\tbar: baz")
	_, err := csyaml.Merge([][]byte{src})
	cstest.RequireErrorContains(t, err, "decoding document 0: [2:1] found character '\t' that cannot start any token")
}

func TestNestedDeepMergePreservesOrder(t *testing.T) {
	left := `
settings:
  ui:
    theme: light
    toolbar:
      - cut
      - copy
`
	right := `
settings:
  ui:
    toolbar:
    - paste
    security: strict
`
	expect := `settings:
  ui:
    theme: light
    toolbar:
    - paste
    security: strict
`
	merged, err := csyaml.Merge([][]byte{[]byte(left), []byte(right)})
	require.NoError(t, err)
	assert.Equal(t, expect, merged.String())
}

// Don't coerce boolean-like strings to true/false (YAML 1.2 / goccy/go-yaml behavior).
func TestBooleanNoCoercion(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"foo: yes", `foo: "yes"`},
		{"foo: YES", `foo: "YES"`},
		{"foo: no", `foo: "no"`},
		{"foo: NO", `foo: "NO"`},
		{"foo: on", `foo: "on"`},
		{"foo: ON", `foo: "ON"`},
		{"foo: off", `foo: "off"`},
		{"foo: OFF", `foo: "OFF"`},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			buf, err := csyaml.Merge([][]byte{nil, []byte(tt.in)})
			require.NoError(t, err)
			assert.Equal(t, tt.out, strings.TrimSuffix(buf.String(), "\n"))
		})
	}
}

// Do coerce boolean-like values to true/false (YAML 1.1 / yaml.v3 behavior).
// func TestBooleanCoercion(t *testing.T) {
//	tests := []struct {
//		in, out string
//	}{
//		{"yes\n", "true\n"},
//		{"YES\n", "true\n"},
//		{"no\n", "false\n"},
//		{"NO\n", "false\n"},
//		{"on\n", "true\n"},
//		{"ON\n", "true\n"},
//		{"off\n", "false\n"},
//		{"OFF\n", "false\n"},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.in, func(t *testing.T) {
//			buf, err := csyaml.Merge([][]byte{nil, []byte(tt.in)})
//			require.NoError(t, err)
//			assert.Equal(t, tt.out, buf.String())
//		})
//	}
//}

func TestExplicitNilOverride(t *testing.T) {
	base := []byte("foo: {one: two}\n")
	override := []byte("foo: ~\n")
	merged, err := csyaml.Merge([][]byte{base, override})
	require.NoError(t, err)
	assert.Equal(t, "foo: null\n", merged.String())
}

func TestOrderPreservation(t *testing.T) {
	left := []byte("a: 1\nb: 2\n")
	right := []byte("c: 3\nb: 20\n")
	expect := "a: 1\nb: 20\nc: 3\n"
	merged, err := csyaml.Merge([][]byte{left, right})
	require.NoError(t, err)
	assert.Equal(t, expect, merged.String())
}
