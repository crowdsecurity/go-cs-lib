package version

import (
	"github.com/blackfireio/osinfo"
)

func DetectOS() (string, string) {
	osInfo, err := osinfo.GetOSInfo()
	if err != nil {
		return System, "???"
	}

	if osInfo.Name != "" && System == "docker" {
		return osInfo.Name + " (docker)", osInfo.Version
	}

	return osInfo.Name, osInfo.Version
}
