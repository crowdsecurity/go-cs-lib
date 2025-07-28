package csyaml

import (
	"bytes"
	"io"
)

// SplitDocuments returns a slice of byte slices, each representing a YAML document.
//
// Since preserving formatting and comments is important but the existing go packages
// all have some issue, this function attempts two strategies: one that decodes and
// re-encodes the YAML content, and another that simply splits the input text.
// If both methods return the same number of documents, we assume the text-based
// function is sufficient. It retains comments and formatting better.
// Otherwise, the round-trip version is used. It retains comments but
// the formatting may be off. The semantics of the document will still be the same
// but if it contains parsing errors, they may refer to a wrong line or column.
//
// This function returns reading errors but any parsing errors are ignored and
// trigger the text-based splitting method.
func SplitDocuments(r io.Reader) ([][]byte, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	textDocs, errText := SplitDocumentsText(bytes.NewReader(input))
	decEncDocs, errDecEnc := SplitDocumentsDecEnc(bytes.NewReader(input))

	if errDecEnc == nil && len(decEncDocs) != len(textDocs) {
		return decEncDocs, nil
	}

	return textDocs, errText
}
