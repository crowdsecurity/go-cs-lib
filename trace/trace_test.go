package trace

import (
	"bytes"
	"os"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestKeeper(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	previousKeeper := keeper
	keeper = &traceKeeper{
		mutex:         &sync.Mutex{},
		dir:           dir,
		crashFileGlob: "crowdsec-crash.*.txt",
		removeAfter:   30 * 24 * time.Hour,
		keepMaxFiles:  100,
	}

	t.Cleanup(func() {
		keeper = previousKeeper
	})
}

func disableLogFatalExit(t *testing.T) {
	t.Helper()

	logger := log.StandardLogger()
	previousExitFunc := logger.ExitFunc
	logger.ExitFunc = func(int) {}

	t.Cleanup(func() {
		logger.ExitFunc = previousExitFunc
	})
}

func captureLogOutput(t *testing.T) *bytes.Buffer {
	t.Helper()

	logger := log.StandardLogger()
	previousOutput := logger.Out
	previousLevel := logger.GetLevel()
	output := &bytes.Buffer{}
	logger.SetOutput(output)
	logger.SetLevel(log.InfoLevel)

	t.Cleanup(func() {
		logger.SetOutput(previousOutput)
		logger.SetLevel(previousLevel)
	})

	return output
}

func TestWriteStackTraceAndList(t *testing.T) {
	setupTestKeeper(t)

	filename, err := WriteStackTrace("write-test-error")
	require.NoError(t, err)

	files, err := List()
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, filename, files[0])

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	assert.Contains(t, string(content), "error: write-test-error")
}

func TestCatchRecoversPanic(t *testing.T) {
	setupTestKeeper(t)
	disableLogFatalExit(t)

	func() {
		defer ReportPanic()
		panic("catch-test-panic")
	}()

	files, err := List()
	require.NoError(t, err)
	require.Len(t, files, 1)

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "error: catch-test-panic")
}

func TestCatchPanicRecoversPanic(t *testing.T) {
	setupTestKeeper(t)
	disableLogFatalExit(t)

	func() {
		defer CatchPanic("some-function")
		panic("catch-panic-test")
	}()

	files, err := List()
	require.NoError(t, err)
	require.Len(t, files, 1)

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "error: catch-panic-test")
}

func TestReportPanicLogsStackTrace(t *testing.T) {
	setupTestKeeper(t)
	disableLogFatalExit(t)
	output := captureLogOutput(t)

	func() {
		defer ReportPanic()
		panic("report-panic-stack")
	}()

	assert.Contains(t, output.String(), "runtime/debug.Stack()")
}
