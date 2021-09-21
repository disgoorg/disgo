package info

import (
	"runtime"
	"runtime/debug"
	"strings"
)

//goland:noinspection GoUnusedConst
const (
	GitHub = "https://github.com/DisgoOrg/disgo"
	Name   = "disgo"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	Version = getVersion()
	OS      = getOS()
)

func getVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if strings.Contains(GitHub, dep.Path) {
				return dep.Version
			}
		}
	}
	return "unknown"
}

func getOS() string {
	os := runtime.GOOS
	if strings.HasPrefix(os, "windows") {
		return "windows"
	}
	if strings.HasPrefix(os, "darwin") {
		return "darwin"
	}
	return "linux"
}
