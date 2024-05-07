package downloader

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// NotFoundError returned by Download() when the remote resource is not found.
type NotFoundError struct {
	URL string
}

func (e NotFoundError) Error() string {
	return "document not found at " + e.URL
}

// BadHTTPCodeError returned when the server responds with unexpected HTTP code.
type BadHTTPCodeError struct {
	URL  string
	Code int
}

func (e BadHTTPCodeError) Error() string {
	return fmt.Sprintf("bad HTTP code %d for %s", e.Code, e.URL)
}

func nullLogger() *logrus.Entry {
	log := logrus.New()
	log.SetOutput(io.Discard)

	return logrus.NewEntry(log)
}

// SHA256 returns the hash of the file if possible, empty string otherwise.
func SHA256(path string) (string, error) {
	file, err := os.Open(path)

	switch {
	case os.IsNotExist(err):
		// first time download
		return "", nil
	case err != nil:
		return "", err
	}

	defer file.Close()

	hash := crypto.SHA256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Downloader fetches a file from a URL to a destination path, with various options.
type Downloader struct {
	// aligned with "betteralign -apply"
	logger             *logrus.Entry
	etagFn             *func(string) (string, error)
	etagPath           string
	httpClient         *http.Client
	destPath           string
	verifyHashFunction string
	verifyHashValue    string
	maxSize            int64
	shelfLife          time.Duration // update if local file is older than this
	mode               os.FileMode
	makeDirs           bool
	ifModifiedSince    bool
	lastModified       bool
	compareContent     bool
}

// New creates a new downloader for the given URL.
func New() *Downloader {
	logger := nullLogger()

	return &Downloader{
		logger:     logger,
		httpClient: http.DefaultClient,
	}
}

// WithLogger sets the logger for the downloader.
// If not set, nothing will be logged.
func (d *Downloader) WithLogger(logger *logrus.Entry) *Downloader {
	d.logger = logger
	return d
}

// ToFile sets the destination path for the downloaded file.
// If a file already exists, its modification time and mode will be used
// by WithLastModified() and WithMode().
func (d *Downloader) ToFile(destPath string) *Downloader {
	d.destPath = destPath
	return d
}

// WithMakeDirs sets whether the downloader should create directories as needed.
// If not set, the downloader will fail if the destination directory does not exist.
func (d *Downloader) WithMakeDirs(makeDirs bool) *Downloader {
	d.makeDirs = makeDirs
	return d
}

// WithETag sets the value to use in the "If-None-Match" header.
// If an ETag is set, a server will not initiate a download
// if the remote file matches the ETag.
func (d *Downloader) WithETag(etag string) *Downloader {
	callback := func(_ string) (string, error) {
		return etag, nil
	}

	// NOTE: we can't check if etagFn is already set, so we can't prevent
	// calling WithETag(), WithETagFn() and WithEtagFile() on the same Downloader.
	// (unless we collect errors along the way and return them in Download()
	d.etagFn = &callback

	return d
}

// WithETagFn sets a function whose result will be used in the "If-None-Match" header.
// The function can compute an ETag from a file's content.
// The alternative is to store it separately (might be a database identifier or uuid, etc)
// and use WithETag() to set it directly.
func (d *Downloader) WithETagFn(etagFn func(string) (string, error)) *Downloader {
	d.etagFn = &etagFn
	return d
}

// IfModifiedSince sets the "If-Modified-Since" header to the file's modification time.
// If the remote resource has not been modified since the given time, the server will
// respond with a 304.
func (d *Downloader) IfModifiedSince() *Downloader {
	d.ifModifiedSince = true
	return d
}

// WithLastModified sets the downloader to check the "Last-Modified" header with a HEAD request.
func (d *Downloader) WithLastModified() *Downloader {
	d.lastModified = true
	return d
}

// WithShelfLife sets the duration after which a file is considered stale, if it has no
// "Last-Modified" header. If unset, the file will be considered stale by default.
func (d *Downloader) WithShelfLife(shelfLife time.Duration) *Downloader {
	d.shelfLife = shelfLife
	return d
}

// WithMode sets the file mode for the downloaded file. If not set, the file mode
// will be taken from the destination file, if it exists before the download.
func (d *Downloader) WithMode(mode os.FileMode) *Downloader {
	d.mode = mode
	return d
}

// WithHTTPClient sets the http client for the downloader.
func (d *Downloader) WithHTTPClient(client *http.Client) *Downloader {
	d.httpClient = client
	return d
}

// LimitDownloadSize sets the maximum size of the downloaded file,
// by checking the Content-Length header and monitoring the size again
// while uncompressing the payload.
func (d *Downloader) LimitDownloadSize(size int64) *Downloader {
	d.maxSize = size
	return d
}

// getDestInfo returns the modification time and file mode of the destination file.
func (d *Downloader) getDestInfo() (time.Time, fs.FileMode) {
	dstInfo, err := os.Stat(d.destPath)

	switch {
	case os.IsNotExist(err):
		return time.Time{}, 0
	case err != nil:
		d.logger.Errorf("Failed to stat destination file %s: %s", d.destPath, err)
		return time.Time{}, 0
	}

	return dstInfo.ModTime(), dstInfo.Mode().Perm()
}

// checkLastModified returns true if the file seems up to date according to its modification time.
func (d *Downloader) checkLastModified(ctx context.Context, url string, modTime time.Time) (bool, error) {
	if !d.lastModified {
		return false, nil
	}

	localIsOld := true

	if d.shelfLife != 0 {
		localIsOld = modTime.Add(d.shelfLife).Before(time.Now())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create HEAD request for %s: %w", url, err)
	}

	client := d.httpClient
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to make HEAD request for %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, BadHTTPCodeError{url, resp.StatusCode}
	}

	remoteLastModified := resp.Header.Get("Last-Modified")
	if remoteLastModified == "" {
		if !localIsOld {
			d.logger.Debugf("No last modified header, but local file is not old: %s",
				d.destPath)
			return true, nil
		}

		d.logger.Debugf("No last modified header: %s", d.destPath)

		return false, nil
	}

	lastAvailable, err := time.Parse(http.TimeFormat, remoteLastModified)
	if err != nil {
		d.logger.Errorf("Failed to parse last modified header '%s': %s", remoteLastModified, err)
	}

	if modTime.After(lastAvailable) {
		d.logger.Debugf("Local file is newer than remote: %s (%s vs %s)",
			d.destPath, modTime, lastAvailable)
		return true, nil
	}

	return false, nil
}

// VerifyHash sets the hash function and value to check the downloaded file against.
func (d *Downloader) VerifyHash(hashFunction, hashValue string) *Downloader {
	d.verifyHashFunction = hashFunction
	d.verifyHashValue = hashValue

	return d
}

func (d *Downloader) selectHashFunction() (hash.Hash, error) {
	switch d.verifyHashFunction {
	case "sha256":
		return crypto.SHA256.New(), nil
	case "md5":
		return crypto.MD5.New(), nil
	case "":
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported hash function %s", d.verifyHashFunction)
	}
}

// ValidateOptions checks that the downloader options are consistent. This is called by Download().
func (d *Downloader) ValidateOptions() error {
	// for the better or worse, due to method chaining we must
	// put all checks in the only place where we can return an error
	if d.destPath == "" {
		return errors.New("destination path must be set")
	}

	if d.shelfLife < 0 {
		return errors.New("shelfLife must not be negative")
	}

	if d.verifyHashFunction != "" && d.verifyHashValue == "" {
		return errors.New("hash value must be set when hash function is set")
	}

	if d.verifyHashFunction == "" && d.verifyHashValue != "" {
		return errors.New("hash function must be set when hash value is set")
	}

	cacheConditions := 0

	if d.lastModified {
		cacheConditions++
	}

	if d.ifModifiedSince {
		cacheConditions++
	}

	if d.etagFn != nil {
		cacheConditions++
	}

	if cacheConditions > 1 {
		return errors.New("only one of lastModified, ifModifiedSince, etagFn can be set")
	}

	return nil
}

// isAtLeastAsRecent() returns true if both files exist and the first one is not older than the second.
// If any of the files does not exist, return false.
func isAtLeastAsRecent(path1, path2 string) bool {
	info1, err := os.Stat(path1)
	if err != nil {
		return false
	}

	info2, err := os.Stat(path2)
	if err != nil {
		return false
	}

	return !info1.ModTime().Before(info2.ModTime())
}

// WithETagFile sets the path to a file where the ETag will be stored and read.
func (d *Downloader) WithETagFile(etagPath string) *Downloader {
	d.etagPath = etagPath
	callback := func(destPath string) (string, error) {
		if !isAtLeastAsRecent(etagPath, destPath) {
			// if the etag file is older than the destination file,
			// it's stale and will be ignored
			return "", nil
		}

		fin, err := os.ReadFile(etagPath)

		switch {
		case os.IsNotExist(err):
			break
		case err != nil:
			return "", fmt.Errorf("can't read etag file %s: %w", etagPath, err)
		}

		return string(fin), nil
	}

	d.etagFn = &callback

	return d
}

// storeETag writes the content of an ETag header to the file.
func storeETag(resp *http.Response, etagPath string, logger *logrus.Entry) {
	if etagPath == "" {
		return
	}

	etag := resp.Header.Get("ETag")
	if etag == "" {
		logger.Warn("No ETag header")
		return
	}

	if err := os.WriteFile(etagPath, []byte(etag), 0o600); err != nil {
		logger.Errorf("Failed to write ETag to %s: %s", etagPath, err)
	}
}

// CompareContent sets the downloader to compare the content after download with
// the previous file at the destination path. If the content is the same, the old
// file will stay in place and be reported as up to date. This is useful when there
// is no other way to check (ETag, Last-Modified...).
func (d *Downloader) CompareContent() *Downloader {
	d.compareContent = true
	return d
}

// compareFiles compares the content of two files and returns true if they are identical.
func compareFiles(file1, file2 string) (bool, error) {
	// NOTE: could be (micro) optimized by comparing sizes first
	f1, err := os.Open(file1)

	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	}

	defer f1.Close()

	f2, err := os.Open(file2)

	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	}

	defer f2.Close()

	const bufSize = 4096
	buf1 := make([]byte, bufSize)
	buf2 := make([]byte, bufSize)

	for {
		n1, err1 := f1.Read(buf1)
		n2, err2 := f2.Read(buf2)

		switch {
		case errors.Is(err1, io.EOF) && errors.Is(err2, io.EOF):
			return true, nil
		case errors.Is(err1, io.EOF) || errors.Is(err2, io.EOF):
			return false, nil
		case err1 != nil || err2 != nil:
			return false, fmt.Errorf("read failed: %w / %w", err1, err2)
		case n1 != n2 || !bytes.Equal(buf1[:n1], buf2[:n2]):
			return false, nil
		}
	}
}

// Download downloads the file from the URL to the destination path.
// Returns true if the file was downloaded, false if it was already up to date.
func (d *Downloader) Download(ctx context.Context, url string) (bool, error) {
	// only one of etagfn, ifmod, lastmod
	if err := d.ValidateOptions(); err != nil {
		return false, fmt.Errorf("downloader options: %w", err)
	}

	d.logger.Debugf("Checking %s", d.destPath)

	destModTime, destFileMode := d.getDestInfo()

	uptodate, err := d.checkLastModified(ctx, url, destModTime)
	if err != nil {
		d.logger.Warnf("Failed to check last modified: %s", err)
	}

	if uptodate {
		return false, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create http request for %s: %w", url, err)
	}

	req.Header.Add("Accept-Encoding", "gzip")

	if d.ifModifiedSince && (destModTime != time.Time{}) {
		req.Header.Add("If-Modified-Since", destModTime.Format(http.TimeFormat))
	}

	// only add If-None-Match if destPath exists,
	// it could have been deleted leaving an .etag
	if d.etagFn != nil && (destModTime != time.Time{}) {
		etag, err := (*d.etagFn)(d.destPath)
		if err != nil {
			d.logger.Warnf("Failed to get etag: %s", err)
		}

		if etag != "" {
			req.Header.Add("If-None-Match", etag)
		}
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed http request for %s: %w", url, err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return false, NotFoundError{url}
	case http.StatusOK:
		break
	case http.StatusNotModified:
		d.logger.Debug("Not modified")
		return false, nil
	default:
		return false, BadHTTPCodeError{url, resp.StatusCode}
	}

	if d.maxSize > 0 {
		contentLengthStr := resp.Header.Get("Content-Length")
		if contentLengthStr != "" {
			contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
			if err != nil {
				d.logger.Warnf("Failed to parse Content-Length header: %s", err)
			} else if contentLength > d.maxSize {
				return false, fmt.Errorf(
					"refusing to download file larger than %d bytes: Content-Length=%d",
					d.maxSize, contentLength)
			}
		}
	}

	reader := resp.Body

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(reader)
		if err != nil {
			return false, fmt.Errorf("failed to create gzip reader: %w", err)
		}

		defer gzipReader.Close()
		reader = gzipReader
	}

	if d.maxSize > 0 {
		reader = NewLimitedReader(reader, d.maxSize)
	}

	destDir, destName := filepath.Split(d.destPath)

	if d.makeDirs {
		if err = os.MkdirAll(destDir, 0o755); err != nil {
			return false, fmt.Errorf("failed to create directories for %s: %w", d.destPath, err)
		}
	}

	tmpFile, err := os.CreateTemp(destDir, destName+".*.download")
	if err != nil {
		return false, fmt.Errorf("failed to create temporary download file for %s: %w", d.destPath, err)
	}

	tmpFileName := tmpFile.Name()
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFileName)
	}()

	// update the file mode from the options, or the pre-existing file mode, if any

	fileMode := d.mode
	if fileMode == 0 && destFileMode != 0 {
		fileMode = destFileMode
	}

	if fileMode != 0 {
		if err = tmpFile.Chmod(fileMode); err != nil {
			return false, fmt.Errorf("failed to chmod temporary file %s: %w", d.destPath, err)
		}
	}

	hasher, err := d.selectHashFunction()
	if err != nil {
		return false, fmt.Errorf("while hashing %s: %w", d.destPath, err)
	}

	writers := []io.Writer{tmpFile}
	if hasher != nil {
		writers = append(writers, hasher)
	}

	multiWriter := io.MultiWriter(writers...)

	written, err := io.Copy(multiWriter, reader)

	switch {
	case errors.Is(err, ErrSizeLimitExceeded):
		return false, fmt.Errorf("download of %s halted: limit of %d bytes exceeded", tmpFileName, d.maxSize)
	case err != nil:
		return false, fmt.Errorf("while writing to %s: %w", tmpFileName, err)
	}

	d.logger.Debugf("Written %d bytes to %s", written, d.destPath)

	if hasher != nil {
		got := hex.EncodeToString(hasher.Sum(nil))
		if got != d.verifyHashValue {
			return false, fmt.Errorf("hash mismatch: expected %s, got %s", d.verifyHashValue, got)
		}
	}

	if err = tmpFile.Sync(); err != nil {
		return false, err
	}

	if err = tmpFile.Close(); err != nil {
		return false, err
	}

	storeETag(resp, d.etagPath, d.logger)

	if d.compareContent {
		same, err := compareFiles(d.destPath, tmpFileName)
		if err != nil {
			d.logger.Errorf("Failed to compare files, assuming different: %s", err)
		}

		if same {
			d.logger.Debugf("Content is the same, not replacing %s", d.destPath)
			// still, we need to update the modification time
			now := time.Now()
			if err = os.Chtimes(d.destPath, now, now); err != nil {
				return false, err
			}

			return false, nil
		}
	}

	if runtime.GOOS == "windows" {
		// On Windows, rename will fail if the destination file already exists
		// so we remove it first.
		err = os.Remove(d.destPath)

		switch {
		case errors.Is(err, fs.ErrNotExist):
			break
		case err != nil:
			d.logger.Errorf("Failed to remove destination file before renaming: %s", err)
			return false, err
		}
	}

	if err = os.Rename(tmpFileName, d.destPath); err != nil {
		return false, err
	}

	return true, nil
}
