package csdaemon

import (
	"os"

	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/sirupsen/logrus"
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

// Notify systemd that the service is ready.
func NotifySystemd(log logrus.FieldLogger) error {
	sent, err := daemon.SdNotify(false, daemon.SdNotifyReady)
	if sent {
		log.Debug("systemd notified")
		return err
	}
	if !DetectSystemd() {
		log.Debug("not running under systemd")
		return nil
	}
	if err != nil {
		log.Error("Failed to notify systemd: %w", err)
		return err
	}
	log.Warn("Systemd notification not supported")
	return nil
}
