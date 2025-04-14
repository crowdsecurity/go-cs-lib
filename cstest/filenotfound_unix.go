//go:build unix || linux || freebsd || netbsd || openbsd || solaris

package cstest

const (
	// these are the same on unix.
	FileNotFoundMessage = "no such file or directory"
	PathNotFoundMessage = FileNotFoundMessage
)
