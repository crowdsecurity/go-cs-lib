package csyaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// SplitDocumentsDecEnc splits documents from reader and returns them as
// re-encoded []byte slices, preserving comments but not exact original
// whitespace. It returns an error if any document cannot be decoded or
// re-encoded.
func SplitDocumentsDecEnc(r io.Reader) ([][]byte, error) {
	dec := yaml.NewDecoder(r)

	var docs [][]byte

	idx := 0

	for {
		var node yaml.Node
		if err := dec.Decode(&node); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("decode doc %d: %w", idx, err)
		}

		var buf bytes.Buffer

		enc := yaml.NewEncoder(&buf)
		enc.SetIndent(2)
		if err := enc.Encode(&node); err != nil {
			return nil, fmt.Errorf("encode doc %d: %w", idx, err)
		}

		_ = enc.Close()

		docs = append(docs, buf.Bytes())
		idx++
	}

	return docs, nil
}
