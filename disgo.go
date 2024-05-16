// Package disgo is a collection of packages for interaction with the Discord Bot and OAuth2 API.
//
// # Discord
//
// Package discord is a collection of structs and types of the Discord API.
//
// # Bot
//
// Package bot connects the Gateway/Sharding, HTTPServer, Cache, Rest & Events packages into a single high level client interface.
//
// # Gateway
//
// Package gateway is used to connect and interact with the Discord Gateway.
//
// # Sharding
//
// Package sharding is used to connect and interact with the Discord Gateway.
//
// # Cache
//
// Package cache provides a generic cache interface for Discord entities.
//
// # HTTPServer
//
// Package httpserver is used to interact with the Discord outgoing webhooks for interactions.
//
// # Events
//
// Package events provide high level events around the Discord Events.
//
// # Rest
//
// Package rest is used to interact with the Discord REST API.
//
// # Webhook
//
// Package webhook provides a high level client interface for interacting with Discord webhooks.
//
// # OAuth2
//
// Package oauth2 provides a high level client interface for interacting with Discord oauth2.
//
// # Voice
//
// Package voice provides a high level client interface for interacting with Discord voice.
package disgo

import (
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/handlers"
)

const (
	// Name is the library name
	Name = "disgo"
	// Module is the library module name
	Module = "github.com/disgoorg/disgo"
	// GitHub is a link to the libraries GitHub repository
	GitHub = "https://github.com/disgoorg/disgo"
)

var (
	// Version is the currently used version of DisGo
	Version = getVersion()

	SemVersion = "semver:" + Version

	// OS is the currently used OS
	OS = getOS()
)

func getVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if dep.Path == Module {
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

// New creates a new bot.Client with the provided token & bot.ConfigOpt(s)
func New(token string, opts ...bot.ConfigOpt) (bot.Client, error) {
	config := bot.DefaultConfig(handlers.GetGatewayHandlers(), handlers.GetHTTPServerHandler())
	config.Apply(opts)

	return bot.BuildClient(token,
		config,
		handlers.DefaultGatewayEventHandlerFunc,
		handlers.DefaultHTTPServerEventHandlerFunc,
		OS,
		Name,
		GitHub,
		Version,
	)
}
