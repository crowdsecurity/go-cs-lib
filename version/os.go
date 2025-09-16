package version

import (
	"github.com/shirou/gopsutil/v4/host"
)

func DetectOS() (string, string, string) {
	info, err := host.Info()
	if err != nil {
		return System, "???", "???"
	}

	platform := info.Platform
	family := info.PlatformFamily
	version := info.PlatformVersion

	if platform != "" && System == "docker" {
		return platform + " (docker)", family, version
	}

	return platform, family, version
}
