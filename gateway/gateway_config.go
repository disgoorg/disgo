package gateway

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway/grate"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
)

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	LargeThreshold:    50,
	GatewayIntents:    discord.GatewayIntentsDefault,
	Compress:          true,
	AutoReconnect:     true,
	MaxReconnectTries: 10,
	OS:                info.OS,
	Browser:           info.Name,
	Device:            info.Name,
}

type Config struct {
	Logger               log.Logger
	EventHandlerFunc     EventHandlerFunc
	LargeThreshold       int
	GatewayIntents       discord.GatewayIntents
	Compress             bool
	SessionID            *string
	LastSequenceReceived *discord.GatewaySequence
	AutoReconnect        bool
	MaxReconnectTries    int
	RateLimiter          grate.Limiter
	RateLimiterConfig    *grate.Config
	Presence             *discord.UpdatePresenceCommandData
	OS                   string
	Browser              string
	Device               string
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithLargeThreshold(largeThreshold int) ConfigOpt {
	return func(config *Config) {
		config.LargeThreshold = largeThreshold
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayIntents(gatewayIntents ...discord.GatewayIntents) ConfigOpt {
	return func(config *Config) {
		var intents discord.GatewayIntents
		for _, intent := range gatewayIntents {
			intents = intents.Add(intent)
		}
		config.GatewayIntents = intents
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCompress(compress bool) ConfigOpt {
	return func(config *Config) {
		config.Compress = compress
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithSessionID(sessionID string) ConfigOpt {
	return func(config *Config) {
		config.SessionID = &sessionID
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithSequence(sequence discord.GatewaySequence) ConfigOpt {
	return func(config *Config) {
		config.LastSequenceReceived = &sequence
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithAutoReconnect(autoReconnect bool) ConfigOpt {
	return func(config *Config) {
		config.AutoReconnect = autoReconnect
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithMaxReconnectTries(maxReconnectTries int) ConfigOpt {
	return func(config *Config) {
		config.MaxReconnectTries = maxReconnectTries
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRateLimiter(rateLimiter grate.Limiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRateLimiterConfig(rateLimiterConfig grate.Config) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfig = &rateLimiterConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRateLimiterConfigOpts(opts ...grate.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		if config.RateLimiterConfig == nil {
			config.RateLimiterConfig = &grate.DefaultConfig
		}
		config.RateLimiterConfig.Apply(opts)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithPresence(presence discord.UpdatePresenceCommandData) ConfigOpt {
	return func(config *Config) {
		config.Presence = &presence
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithOS(os string) ConfigOpt {
	return func(config *Config) {
		config.OS = os
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithBrowser(browser string) ConfigOpt {
	return func(config *Config) {
		config.Browser = browser
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithDevice(device string) ConfigOpt {
	return func(config *Config) {
		config.Device = device
	}
}
