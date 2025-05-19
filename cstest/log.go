package cstest

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
)

// CaptureLogs captures logs from the standard logger and returns a buffer.
func CaptureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()

	// grab the old state
	std := logrus.StandardLogger()
	oldOut := std.Out
	oldFmt := std.Formatter
	oldLvl := std.GetLevel()

	// set up our buffer
	buf := &bytes.Buffer{}
	std.SetOutput(buf)
	std.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, DisableColors: true})
	std.SetLevel(logrus.DebugLevel)

	// restore on cleanup
	t.Cleanup(func() {
		std.SetOutput(oldOut)
		std.SetFormatter(oldFmt)
		std.SetLevel(oldLvl)
	})

	return buf
}
