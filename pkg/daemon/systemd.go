package daemon

import (
	"os"
)

// DetectSystemd detects if the current process is running under systemd
// and returns true if it is.
// It is not intended to be 100% reliable, but useful to guess if we
// should treat a notification failure as an error.
func DetectSystemd() bool {
	if os.Getenv("INVOCATION_ID") != "" {
		return true
	}
	if os.Getenv("JOURNAL_STREAM") != "" {
		return true
	}
	if os.Getenv("NOTIFY_SOCKET") != "" {
		return true
	}
	return false
}
