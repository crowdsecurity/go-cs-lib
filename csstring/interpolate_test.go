package csstring_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crowdsecurity/go-cs-lib/csstring"
	"github.com/crowdsecurity/go-cs-lib/cstest"
)

type person struct {
	Name string
}

func TestInterpolate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		data        any
		expected    string
		expectedErr string
	}{
		{
			name:     "Successful map interpolation",
			template: "{{.Name}}",
			data:     map[string]string{"Name": "JohnDoe"},
			expected: "JohnDoe",
		},
		{
			name:     "Successful struct interpolation",
			template: "{{.Name}}",
			data:     person{Name: "JaneDoe"},
			expected: "JaneDoe",
		},
		{
			name:        "Unsuccessful interpolation, missing key",
			template:    "{{.Name}}",
			data:        map[string]string{"FirstName": "JohnDoe"},
			expectedErr: "template: :1:2: executing \"\" at <.Name>: map has no entry for key \"Name\"",
		},
		{
			name:        "Unsuccessful interpolation, missing field",
			template:    "{{.FullName}}",
			data:        person{Name: "JohnDoe"},
			expectedErr: "template: :1:2: executing \"\" at <.FullName>: can't evaluate field FullName in type csstring_test.person",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := csstring.Interpolate(test.template, test.data)
			cstest.RequireErrorContains(t, err, test.expectedErr)
			assert.Equal(t, test.expected, res)
		})
	}
}
