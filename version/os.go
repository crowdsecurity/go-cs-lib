package version

import (
	"github.com/shirou/gopsutil/v4/host"
)

func DetectOS() (string, string) {
	info, err := host.Info()
	if err != nil {
		return System, "???"
	}

	name := info.Platform
	version := info.PlatformVersion

	if name != "" && System == "docker" {
		return name + " (docker)", version
	}

	return name, version
}
