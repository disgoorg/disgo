package disgo

import (
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/gateway/handlers"
	"github.com/disgoorg/disgo/httpserver"
)

//goland:noinspection GoUnusedConst
const (
	Name   = "disgo"
	GitHub = "https://github.com/disgoorg/" + Name
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

// New creates a new core.Client instance with the provided bot token & ConfigOpt(s)
//goland:noinspection GoUnusedExportedFunction
func New(token string, opts ...bot.ConfigOpt) (bot.Client, error) {
	config := bot.DefaultConfig(handlers.GetGatewayHandlers(), handlers.GetHTTPServerHandler())
	config.Apply(opts)

	return bot.BuildClient(token, *config,
		func(client bot.Client) gateway.EventHandlerFunc {
			return handlers.DefaultGatewayEventHandler(client)
		},
		func(client bot.Client) httpserver.EventHandlerFunc {
			return handlers.DefaultHTTPServerEventHandler(client)
		},
		OS, Name, GitHub, Version,
	)
}
