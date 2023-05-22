package version

import (
	"fmt"
	"log"
	"runtime"
)

var (
	Version   string                  // = "v0.0.0"
	BuildDate string                  // = "2023-03-06_09:55:34"
	Tag       string                  // = "dev"
	GoVersion = runtime.Version()[2:] // = "1.13"
)

func ShowStr() string {
	ret := ""
	ret += fmt.Sprintf("version: %s-%s\n", Version, Tag)
	ret += fmt.Sprintf("BuildDate: %s\n", BuildDate)
	ret += fmt.Sprintf("GoVersion: %s\n", GoVersion)
	return ret
}

func Show() {
	log.Printf("version: %s-%s", Version, Tag)
	log.Printf("BuildDate: %s", BuildDate)
	log.Printf("GoVersion: %s", GoVersion)
}

func VersionStr() string {
	return fmt.Sprintf("%s-%s", Version, Tag)
}
