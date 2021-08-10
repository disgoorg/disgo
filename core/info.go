package core

import (
	"runtime/debug"
	"strings"
)

// GitHub is the Disgo GitHub URL
const GitHub = "https://github.com/DisgoOrg/disgo"

// Version returns the current used Disgo version in the format vx.x.x
var Version = getVersion()

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
