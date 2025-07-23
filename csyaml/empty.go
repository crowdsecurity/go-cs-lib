package csyaml

import (
	"errors"
	"io"

	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

// IsEmptyYAML reads one or more YAML documents from r and returns true
// if they are all empty or contain only comments.
// It will reports errors if the input is not valid YAML.
func IsEmptyYAML(r io.Reader) (bool, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return false, err
	}

	if len(src) == 0 {
		return true, nil
	}

	file, err := parser.ParseBytes(src, 0)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return true, nil
		}

		return false, errors.New(yaml.FormatError(err, false, false))
	}

	if file == nil || len(file.Docs) == 0 {
		return true, nil
	}

	for _, doc := range file.Docs {
		if doc.Body != nil {
			return false, nil
		}
	}

	return true, nil
}
