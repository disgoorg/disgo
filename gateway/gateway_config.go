package gateway

import (
	"github.com/disgoorg/log"
	"github.com/gorilla/websocket"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Logger:            log.Default(),
		Dialer:            websocket.DefaultDialer,
		LargeThreshold:    50,
		Intents:           IntentsDefault,
		Compress:          true,
		ShardID:           0,
		ShardCount:        1,
		AutoReconnect:     true,
		MaxReconnectTries: 10,
	}
}

// Config lets you configure your Gateway instance.
type Config struct {
	Logger                    log.Logger
	Dialer                    *websocket.Dialer
	LargeThreshold            int
	Intents                   Intents
	Compress                  bool
	URL                       string
	ShardID                   int
	ShardCount                int
	SessionID                 *string
	LastSequenceReceived      *int
	AutoReconnect             bool
	MaxReconnectTries         int
	EnableRawEvents           bool
	RateLimiter               RateLimiter
	RateRateLimiterConfigOpts []RateLimiterConfigOpt
	Presence                  *MessageDataPresenceUpdate
	OS                        string
	Browser                   string
	Device                    string
}

// ConfigOpt is a type alias for a function that takes a Config and is used to configure your Server.
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	if c.RateLimiter == nil {
		c.RateLimiter = NewRateLimiter(c.RateRateLimiterConfigOpts...)
	}
}

// WithLogger sets the Logger for the Gateway.
func WithLogger(logger log.Logger) ConfigOpt {
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

// WithMaxReconnectTries sets the maximum number of reconnect attempts before stopping.
func WithMaxReconnectTries(maxReconnectTries int) ConfigOpt {
	return func(config *Config) {
		config.MaxReconnectTries = maxReconnectTries
	}
}

// WithEnableRawEvents enables/disables the EventTypeRaw.
func WithEnableRawEvents(enableRawEventEvents bool) ConfigOpt {
	return func(config *Config) {
		config.EnableRawEvents = enableRawEventEvents
	}
}

// WithRateLimiter sets the grate.RateLimiter for the Gateway.
func WithRateLimiter(rateLimiter RateLimiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

// WithRateRateLimiterConfigOpts lets you configure the default RateLimiter.
func WithRateRateLimiterConfigOpts(opts ...RateLimiterConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RateRateLimiterConfigOpts = append(config.RateRateLimiterConfigOpts, opts...)
	}
}

// WithPresence sets the initial presence the bot should display.
func WithPresence(presence MessageDataPresenceUpdate) ConfigOpt {
	return func(config *Config) {
		config.Presence = &presence
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
