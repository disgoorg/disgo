package gateway

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Logger:          slog.Default(),
		Dialer:          websocket.DefaultDialer,
		LargeThreshold:  50,
		Intents:         IntentsDefault,
		Compress:        true,
		URL:             "wss://gateway.discord.gg",
		ShardID:         0,
		ShardCount:      1,
		AutoReconnect:   true,
		EnableResumeURL: true,
	}
}

// Config lets you configure your Gateway instance.
type Config struct {
	// Logger is the Logger of the Gateway. Defaults to slog.Default().
	Logger *slog.Logger
	// Dialer is the websocket.Dialer of the Gateway. Defaults to websocket.DefaultDialer.
	Dialer *websocket.Dialer
	// LargeThreshold is the threshold for the Gateway. Defaults to 50
	// See here for more information: https://discord.com/developers/docs/topics/gateway-events#identify-identify-structure.
	LargeThreshold int
	// Intents is the Intents for the Gateway. Defaults to IntentsNone.
	Intents Intents
	// Compress is whether the Gateway should compress payloads. Defaults to true.
	Compress bool
	// URL is the URL of the Gateway. Defaults to fetch from Discord.
	URL string
	// ShardID is the shardID of the Gateway. Defaults to 0.
	ShardID int
	// ShardCount is the shardCount of the Gateway. Defaults to 1.
	ShardCount int
	// SessionID is the last sessionID of the Gateway. Defaults to nil (no resume).
	SessionID *string
	// ResumeURL is the last resumeURL of the Gateway. Defaults to nil (no resume).
	ResumeURL *string
	// LastSequenceReceived is the last sequence received by the Gateway. Defaults to nil (no resume).
	LastSequenceReceived *int
	// AutoReconnect is whether the Gateway should automatically reconnect or call the CloseHandlerFunc. Defaults to true.
	AutoReconnect bool
	// EnableRawEvents is whether the Gateway should emit EventRaw. Defaults to false.
	EnableRawEvents bool
	// EnableResumeURL is whether the Gateway should enable the resumeURL. Defaults to true.
	EnableResumeURL bool
	// RateLimiter is the RateLimiter of the Gateway. Defaults to NewRateLimiter().
	RateLimiter RateLimiter
	// RateLimiterConfigOpts is the RateLimiterConfigOpts of the Gateway. Defaults to nil.
	RateLimiterConfigOpts []RateLimiterConfigOpt
	// Presence is the presence it should send on login. Defaults to nil.
	Presence *MessageDataPresenceUpdate
	// OS is the OS it should send on login. Defaults to runtime.GOOS.
	OS string
	// Browser is the Browser it should send on login. Defaults to "disgo".
	Browser string
	// Device is the Device it should send on login. Defaults to "disgo".
	Device string
}

// ConfigOpt is a type alias for a function that takes a Config and is used to configure your Server.
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.RateLimiter == nil {
		c.RateLimiter = NewRateLimiter(c.RateLimiterConfigOpts...)
	}
}

// WithLogger sets the Logger for the Gateway.
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithDialer sets the websocket.Dialer for the Gateway.
func WithDialer(dialer *websocket.Dialer) ConfigOpt {
	return func(config *Config) {
		config.Dialer = dialer
	}
}

// WithLargeThreshold sets the threshold for the Gateway.
// See here for more information: https://discord.com/developers/docs/topics/gateway#identify-identify-structure
func WithLargeThreshold(largeThreshold int) ConfigOpt {
	return func(config *Config) {
		config.LargeThreshold = largeThreshold
	}
}

// WithIntents sets the Intents for the Gateway.
// See here for more information: https://discord.com/developers/docs/topics/gateway#gateway-intents
func WithIntents(intents ...Intents) ConfigOpt {
	return func(config *Config) {
		config.Intents = config.Intents.Add(intents...)
	}
}

// WithCompress sets whether this Gateway supports compression.
// See here for more information: https://discord.com/developers/docs/topics/gateway#encoding-and-compression
func WithCompress(compress bool) ConfigOpt {
	return func(config *Config) {
		config.Compress = compress
	}
}

// WithURL sets the Gateway URL for the Gateway.
func WithURL(url string) ConfigOpt {
	return func(config *Config) {
		config.URL = url
	}
}

// WithShardID sets the shard ID for the Gateway.
// See here for more information on sharding: https://discord.com/developers/docs/topics/gateway#sharding
func WithShardID(shardID int) ConfigOpt {
	return func(config *Config) {
		config.ShardID = shardID
	}
}

// WithShardCount sets the shard count for the Gateway.
// See here for more information on sharding: https://discord.com/developers/docs/topics/gateway#sharding
func WithShardCount(shardCount int) ConfigOpt {
	return func(config *Config) {
		config.ShardCount = shardCount
	}
}

// WithSessionID sets the Session ID for the Gateway.
// If sessionID and lastSequence is present while connecting, the Gateway will try to resume the session.
func WithSessionID(sessionID string) ConfigOpt {
	return func(config *Config) {
		config.SessionID = &sessionID
	}
}

// WithSequence sets the last sequence received for the Gateway.
// If sessionID and lastSequence is present while connecting, the Gateway will try to resume the session.
func WithSequence(sequence int) ConfigOpt {
	return func(config *Config) {
		config.LastSequenceReceived = &sequence
	}
}

// WithAutoReconnect sets whether the Gateway should automatically reconnect to Discord.
func WithAutoReconnect(autoReconnect bool) ConfigOpt {
	return func(config *Config) {
		config.AutoReconnect = autoReconnect
	}
}

// WithEnableRawEvents enables/disables the EventTypeRaw.
func WithEnableRawEvents(enableRawEventEvents bool) ConfigOpt {
	return func(config *Config) {
		config.EnableRawEvents = enableRawEventEvents
	}
}

// WithEnableResumeURL enables/disables usage of resume URLs sent by Discord.
func WithEnableResumeURL(enableResumeURL bool) ConfigOpt {
	return func(config *Config) {
		config.EnableResumeURL = enableResumeURL
	}
}

// WithRateLimiter sets the grate.RateLimiter for the Gateway.
func WithRateLimiter(rateLimiter RateLimiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

// WithRateLimiterConfigOpts lets you configure the default RateLimiter.
func WithRateLimiterConfigOpts(opts ...RateLimiterConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterConfigOpts = append(config.RateLimiterConfigOpts, opts...)
	}
}

// WithPresenceOpts allows to pass initial presence data the bot should display.
func WithPresenceOpts(opts ...PresenceOpt) ConfigOpt {
	return func(config *Config) {
		presenceUpdate := &MessageDataPresenceUpdate{}
		for _, opt := range opts {
			opt(presenceUpdate)
		}
		config.Presence = presenceUpdate
	}
}

// WithOS sets the operating system the bot is running on.
// See here for more information: https://discord.com/developers/docs/topics/gateway#identify-identify-connection-properties
func WithOS(os string) ConfigOpt {
	return func(config *Config) {
		config.OS = os
	}
}

// WithBrowser sets the browser the bot is running on.
// See here for more information: https://discord.com/developers/docs/topics/gateway#identify-identify-connection-properties
func WithBrowser(browser string) ConfigOpt {
	return func(config *Config) {
		config.Browser = browser
	}
}

// WithDevice sets the device the bot is running on.
// See here for more information: https://discord.com/developers/docs/topics/gateway#identify-identify-connection-properties
func WithDevice(device string) ConfigOpt {
	return func(config *Config) {
		config.Device = device
	}
}
