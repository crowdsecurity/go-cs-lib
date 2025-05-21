package csyaml

import (
	"errors"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
)

// SplitDocuments reads every YAML document from r and returns each
// as a separate []byte.
//
// Documents are round-tripped through goccy’s Decoder/Marshal pipeline,
// comments and exact original formatting may be lost.
func SplitDocuments(r io.Reader) ([][]byte, error) {
	docs := make([][]byte, 0)

	dec := yaml.NewDecoder(r)

	idx := -1

	for {
		var v any

		idx++

		if err := dec.Decode(&v); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decoding document %d: %s", idx, yaml.FormatError(err, false, false))
		}

		out, err := yaml.Marshal(v)
		if err != nil {
			return nil, err
		}

		docs = append(docs, out)
	}

	return docs, nil
}
