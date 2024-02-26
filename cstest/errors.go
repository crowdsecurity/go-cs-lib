package cstest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	logtest "github.com/sirupsen/logrus/hooks/test"
)

// The following functions are used to test for errors and log messages.
// It would be nice to rely less on the content of error messages, and use a more
// structured approach with error types, but that would require both test and code refactoring.

func AssertErrorContains(t *testing.T, err error, expectedErr string) {
	t.Helper()

	if expectedErr != "" {
		assert.ErrorContains(t, err, expectedErr)
		return
	}

	assert.NoError(t, err)
}

func AssertErrorMessage(t *testing.T, err error, expectedErr string) {
	t.Helper()

	if expectedErr != "" {
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}

		assert.Equal(t, expectedErr, errmsg)

		return
	}

	require.NoError(t, err)
}

func RequireErrorContains(t *testing.T, err error, expectedErr string) {
	t.Helper()

	if expectedErr != "" {
		require.ErrorContains(t, err, expectedErr)
		return
	}

	require.NoError(t, err)
}

func RequireErrorMessage(t *testing.T, err error, expectedErr string) {
	t.Helper()

	if expectedErr != "" {
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}

		require.Equal(t, expectedErr, errmsg)

		return
	}

	require.NoError(t, err)
}

func RequireLogContains(t *testing.T, hook *logtest.Hook, expected string) {
	t.Helper()

	// look for a log entry that matches the expected message
	for _, entry := range hook.AllEntries() {
		if strings.Contains(entry.Message, expected) {
			return
		}
	}

	// show all hook entries, in case the test fails we'll need them
	for _, entry := range hook.AllEntries() {
		t.Logf("log entry: %s", entry.Message)
	}

	require.Fail(t, "no log entry found with message", expected)
}
