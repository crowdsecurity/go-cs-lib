package downloader

import (
	"errors"
	"io"
)

var ErrSizeLimitExceeded = errors.New("size limit exceeded")

// NewLimitedReader wraps a reader to read up to n bytes before stopping with ErrSizeLimitExceeded.
// The underlying reader is not closed.
func NewLimitedReader(r io.Reader, n int64) io.ReadCloser {
	return &limitedReader{r: r, n: n}
}

type limitedReader struct {
	r io.Reader
	n int64
}

func (l *limitedReader) Read(p []byte) (int, error) {
	if l.n <= 0 {
		return 0, ErrSizeLimitExceeded
	}

	if int64(len(p)) > l.n {
		p = p[0:l.n]
	}

	n, err := l.r.Read(p)
	l.n -= int64(n)

	return n, err
}

func (*limitedReader) Close() error {
	// closing the underlying reader is left to the caller
	return nil
}
