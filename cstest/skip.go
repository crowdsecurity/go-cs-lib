package cstest

import (
	"os"
	"runtime"
	"testing"
)

func SkipOnWindows(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on windows")
	}
}

func SkipOnWindowsBecause(t *testing.T, reason string) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skipf("Skipping test on windows (%s)", reason)
	}
}

func SkipIfDefined(t *testing.T, envVar string) {
	t.Helper()

	if os.Getenv(envVar) != "" {
		t.Skipf("Skipping test because %s is defined", envVar)
	}
}
