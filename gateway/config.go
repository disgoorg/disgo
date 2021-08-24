package gateway

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/log"
)

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	LargeThreshold: 50,
	GatewayIntents: discord.GatewayIntentsDefault,
	OS:             info.OS,
	Browser:        info.Name,
	Device:         info.Name,
}

type Config struct {
	Logger           log.Logger
	RestServices     rest.Services
	Token            string
	EventHandlerFunc EventHandlerFunc
	LargeThreshold   int
	GatewayIntents   discord.GatewayIntents
	OS               string
	Browser          string
	Device           string
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithLargeThreshold(largeThreshold int) ConfigOpt {
	return func(config *Config) {
		config.LargeThreshold = largeThreshold
	}
}

func WithGatewayIntents(gatewayIntents ...discord.GatewayIntents) ConfigOpt {
	return func(config *Config) {
		var intents discord.GatewayIntents
		for _, intent := range gatewayIntents {
			intents = intents.Add(intent)
		}
		config.GatewayIntents = intents
	}
}

func WithOS(os string) ConfigOpt {
	return func(config *Config) {
		config.OS = os
	}
}

func WithBrowser(browser string) ConfigOpt {
	return func(config *Config) {
		config.Browser = browser
	}
}

func WithDevice(device string) ConfigOpt {
	return func(config *Config) {
		config.Device = device
	}
}
