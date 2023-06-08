package csdaemon

import (
	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/sirupsen/logrus"
)

// Notify systemd that the service is ready.
func NotifySystemd(log logrus.FieldLogger) error {
	sent, err := daemon.SdNotify(false, daemon.SdNotifyReady)
	if sent {
		log.Debug("Systemd notified")
		return err
	}
	if err != nil {
		log.Error("Failed to notify systemd: %w", err)
		return err
	}
	log.Debug("Not running under systemd")
	return nil
}
