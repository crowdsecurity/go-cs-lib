package version

import (
	"fmt"
	"runtime"
)

var (
	Version   string                  // = "v0.0.0"
	BuildDate string                  // = "2023-03-06_09:55:34"
	System    = runtime.GOOS          // = "linux", "windows", "docker" (when overridden by a Makefile)
	Tag       string                  // = "dev"
	GoVersion = runtime.Version()[2:] // = "1.13"
)

func FullString() string {
	ret := fmt.Sprintf("version: %s-%s\n", Version, Tag)
	ret += fmt.Sprintf("BuildDate: %s\n", BuildDate)
	ret += fmt.Sprintf("GoVersion: %s\n", GoVersion)
	ret += fmt.Sprintf("Platform: %s\n", System)

	return ret
}

func String() string {
	return fmt.Sprintf("%s-%s", Version, Tag)
}
