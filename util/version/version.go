package version

import (
	"runtime/debug"
)

// v holds the version number.
var v string

func revision() string {
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				return setting.Value[:7]
			}
		}
	}
	return ""
}

func Version(version string) string {
	if version != "" {
		v = version
	}
	if v == "" {
		v = "dev-" + revision()
	}
	return v
}
