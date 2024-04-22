package downloader_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/crowdsecurity/go-cs-lib/downloader"

	"github.com/stretchr/testify/require"
)

// Download to a temporary location, and compare the downloaded
// content with the destination if it exists. Replace only if changed.

func TestDownloadCompare(t *testing.T) {
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
		ToFile(dest).
		CompareContent()

	// first download, and check the content

	downloaded, err := d.Download(ctx, ts.URL+"/testfile")

	require.True(t, downloaded)
	require.NoError(t, err)

	content, err := os.ReadFile(dest)
	require.Equal(t, "content", string(content))
	require.NoError(t, err)

	// a second download: although a download is happening, the content is the same
	// and the file is not reported as downloaded

	downloaded, err = d.Download(ctx, ts.URL+"/testfile")

	require.False(t, downloaded)
	require.NoError(t, err)

	// now we change the content and redownload

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/testfile":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, "new-content")
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, "not found")
		}
	}))
	defer ts2.Close()

	downloaded, err = d.Download(ctx, ts2.URL+"/testfile")

	require.True(t, downloaded)
	require.NoError(t, err)

	content, err = os.ReadFile(dest)
	require.Equal(t, "new-content", string(content))
	require.NoError(t, err)
}
