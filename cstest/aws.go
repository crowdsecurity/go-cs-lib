package cstest

import (
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

// SetAWSTestEnv sets the environment variables required to run tests against LocalStack,
// including AWS credentials and the custom endpoint.
//
// It also performs basic validation:
//   - Skips the test on Windows or when TEST_LOCAL_ONLY is defined.
//   - Fails the test if AWS_ENDPOINT_FORCE is already set (to avoid unintended overrides).
//   - Fails the test if the LocalStack endpoint is not reachable.
func SetAWSTestEnv(t *testing.T) string {
	t.Helper()

	SkipOnWindows(t)
	SkipIfDefined(t, "TEST_LOCAL_ONLY")

	endpoint := "http://localhost:4566"

	if os.Getenv("AWS_ENDPOINT_FORCE") != "" {
		t.Fatal("AWS_ENDPOINT_FORCE already set -- did you call SetAWSTestEnv() twice?")
	}

	t.Setenv("AWS_ENDPOINT_FORCE", endpoint)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	dialer := &net.Dialer{Timeout: 2 * time.Second}

	_, err := dialer.DialContext(t.Context(), "tcp", strings.TrimPrefix(endpoint, "http://"))
	if err != nil {
		t.Fatalf("%s: make sure localstack is running and retry", err)
	}

	return endpoint
}
