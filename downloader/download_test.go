package downloader_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/crowdsecurity/go-cs-lib/downloader"
)

// simplest case: just download a file every time

func TestDownloadToFile(t *testing.T) {
	ctx := context.Background()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/testfile":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, "content")
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, "not found")
		}
	}))
	defer ts.Close()

	dest := filepath.Join(t.TempDir(), "example.txt")
	defer os.Remove(dest)

	d := downloader.New().
		ToFile(dest)

	// first download, and check the content

	downloaded, err := d.Download(ctx, ts.URL+"/testfile")

	require.True(t, downloaded)
	require.NoError(t, err)

	content, err := os.ReadFile(dest)
	require.Equal(t, "content", string(content))
	require.NoError(t, err)

	// a second download, there are no checks to avoid downloading again

	downloaded, err = d.Download(ctx, ts.URL+"/testfile")

	require.True(t, downloaded)
	require.NoError(t, err)

	// 404

	downloaded, err = d.Download(ctx, ts.URL+"/testfile_missing")

	require.False(t, downloaded)

	var notfound downloader.NotFoundError

	require.ErrorAs(t, err, &notfound, "document not found at "+ts.URL+"/testfile_missing")
}
