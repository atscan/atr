package main

import (
	"github.com/atscan/atr/cmd"
)

var _version string

func main() {
	cmd.Execute()
}

func Version() string {
	return _version
}
