package main

import (
	"github.com/atscan/atr/cmd"
	"github.com/atscan/atr/util/version"
)

var _version string

func main() {
	version.Version(_version)
	cmd.Execute()
}

func Version() string {
	return _version
}
