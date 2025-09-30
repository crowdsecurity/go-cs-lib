//go:build unix || linux || freebsd || netbsd || openbsd || solaris

package cstest

const (
	FileNotFoundMessage = "no such file or directory"
	PathNotFoundMessage = FileNotFoundMessage // these are the same on unix.
)
