package gateway

import (
	"github.com/DisgoOrg/log"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway/grate"
)

//goland:noinspection GoUnusedGlobalVariable
func DefaultConfig() *Config {
	return &Config{
		LargeThreshold:    50,
		GatewayIntents:    discord.GatewayIntentsDefault,
		Compress:          true,
		ShardID:           0,
		ShardCount:        1,
		AutoReconnect:     true,
		MaxReconnectTries: 10,
	}
}

type Config struct {
	Logger                log.Logger
	LargeThreshold        int
	GatewayIntents        discord.GatewayIntents
	Compress              bool
	GatewayURL            string
	ShardID               int
	ShardCount            int
	SessionID             *string
	LastSequenceReceived  *discord.GatewaySequence
	AutoReconnect         bool
	MaxReconnectTries     int
	RateLimiter           grate.Limiter
	RateLimiterConfigOpts []grate.ConfigOpt
	Presence              *discord.UpdatePresenceCommandData
	OS                    string
	Browser               string
	Device                string
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.RateLimiter == nil {
		c.RateLimiter = grate.NewLimiter(c.RateLimiterConfigOpts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
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
func WithGatewayURL(gatewayURL string) ConfigOpt {
	return func(config *Config) {
		config.GatewayURL = gatewayURL
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShardID(shardID int) ConfigOpt {
	return func(config *Config) {
		config.ShardID = shardID
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShardCount(shardCount int) ConfigOpt {
	return func(config *Config) {
		config.ShardCount = shardCount
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
func WithRateLimiterConfigOpts(opts ...grate.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfigOpts = append(config.RateLimiterConfigOpts, opts...)
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
