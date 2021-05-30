package api

import (
	"runtime/debug"
	"strings"
)

const GITHUB = "https://github.com/DisgoOrg/disgo"

var VERSION = getVersion()

func getVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if strings.Contains(GITHUB, dep.Path) {
				return dep.Version
			}
		}
	}
	return "unknown"
}
