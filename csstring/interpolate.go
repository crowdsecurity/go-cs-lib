package csstring

import (
	"strings"
	"text/template"
)

// Interpolate fills a string template with the given values, can be map or struct.
// example: Interpolate("{{.Name}}", map[string]string{"Name": "JohnDoe"})
func Interpolate(s string, data any) (string, error) {
	tmpl, err := template.New("").Option("missingkey=error").Parse(s)
	if err != nil {
		return "", err
	}

	var b strings.Builder

	if err := tmpl.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
