package csyaml

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// SplitDocumentsText splits a YAML input stream into separate documents by looking for the `---` separator.
// No encoding or decoding is performed; the input is treated as raw text.
// Comments and whitespace are preserved. Malformed documents are returned as-is.
func SplitDocumentsText(r io.Reader) ([][]byte, error) {
	var (
		docs    [][]byte
		current bytes.Buffer
	)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		isSeparator := strings.HasPrefix(trimmed, "---") &&
			(trimmed == "---" || strings.HasPrefix(trimmed, "--- "))

		// Always write the line first
		current.WriteString(line)
		current.WriteByte('\n')

		if isSeparator && current.Len() > len(line)+1 { // +1 for newline just added
			// Separator starts a new doc → commit previous one
			// (everything up to this line is the previous doc)
			n := current.Len()
			// rewind to just before this separator line
			doc := current.Bytes()[:n-len(line)-1]
			docs = append(docs, append([]byte(nil), doc...))
			current = *bytes.NewBuffer(current.Bytes()[n-len(line)-1:])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if current.Len() > 0 {
		docs = append(docs, current.Bytes())
	}

	return docs, nil
}
