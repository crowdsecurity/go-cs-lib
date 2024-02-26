package csdaemon

import (
	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/sirupsen/logrus"
)

// Notify systemd that the service is ready.
//
// Deprecated: Use Notify instead.
func NotifySystemd(log logrus.FieldLogger) error {
	return Notify(Ready, log)
}

const (
	Ready = daemon.SdNotifyReady
	Reloading = daemon.SdNotifyReloading
	Stopping = daemon.SdNotifyStopping
	Watchdog = daemon.SdNotifyWatchdog
)

// Notify systemd that the service is ready.
func Notify(state string, log logrus.FieldLogger) error {
	sent, err := daemon.SdNotify(false, state)
	if sent {
		log.Debugf("Systemd notified: %s", state)
		return err
	}
	if err != nil {
		log.Errorf("Failed to notify systemd: %v", err)
		return err
	}
	return nil
}
