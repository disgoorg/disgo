package gateway

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway/grate"
	"github.com/disgoorg/log"
	"github.com/gorilla/websocket"
)

func DefaultConfig() *Config {
	return &Config{
		Logger:            log.Default(),
		Dialer:            websocket.DefaultDialer,
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
	Dialer                *websocket.Dialer
	LargeThreshold        int
	GatewayIntents        discord.GatewayIntents
	Compress              bool
	GatewayURL            string
	ShardID               int
	ShardCount            int
	SessionID             *string
	LastSequenceReceived  *int
	AutoReconnect         bool
	MaxReconnectTries     int
	RateLimiter           grate.Limiter
	RateLimiterConfigOpts []grate.ConfigOpt
	Presence              *discord.GatewayMessageDataPresenceUpdate
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

func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

func WithDialer(dialer *websocket.Dialer) ConfigOpt {
	return func(config *Config) {
		config.Dialer = dialer
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

func WithCompress(compress bool) ConfigOpt {
	return func(config *Config) {
		config.Compress = compress
	}
}

func WithGatewayURL(gatewayURL string) ConfigOpt {
	return func(config *Config) {
		config.GatewayURL = gatewayURL
	}
}

func WithShardID(shardID int) ConfigOpt {
	return func(config *Config) {
		config.ShardID = shardID
	}
}

func WithShardCount(shardCount int) ConfigOpt {
	return func(config *Config) {
		config.ShardCount = shardCount
	}
}

func WithSessionID(sessionID string) ConfigOpt {
	return func(config *Config) {
		config.SessionID = &sessionID
	}
}

func WithSequence(sequence int) ConfigOpt {
	return func(config *Config) {
		config.LastSequenceReceived = &sequence
	}
}

func WithAutoReconnect(autoReconnect bool) ConfigOpt {
	return func(config *Config) {
		config.AutoReconnect = autoReconnect
	}
}

func WithMaxReconnectTries(maxReconnectTries int) ConfigOpt {
	return func(config *Config) {
		config.MaxReconnectTries = maxReconnectTries
	}
}

func WithRateLimiter(rateLimiter grate.Limiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

func WithRateLimiterConfigOpts(opts ...grate.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfigOpts = append(config.RateLimiterConfigOpts, opts...)
	}
}

func WithPresence(presence discord.GatewayMessageDataPresenceUpdate) ConfigOpt {
	return func(config *Config) {
		config.Presence = &presence
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
