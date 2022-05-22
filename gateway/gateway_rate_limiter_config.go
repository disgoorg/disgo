package gateway

import (
	"github.com/disgoorg/log"
)

// DefaultRateLimiterConfig returns a RateLimiterConfig with sensible defaults.
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		Logger:            log.Default(),
		CommandsPerMinute: 120,
	}
}

// RateLimiterConfig lets you configure your Gateway instance.
type RateLimiterConfig struct {
	Logger            log.Logger
	CommandsPerMinute int
}

// RateLimiterConfigOpt is a type alias for a function that takes a RateLimiterConfig and is used to configure your Server.
type RateLimiterConfigOpt func(config *RateLimiterConfig)

// Apply applies the given RateLimiterConfigOpt(s) to the RateLimiterConfig
func (c *RateLimiterConfig) Apply(opts []RateLimiterConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithRateLimiterLogger sets the Logger for the Gateway.
func WithRateLimiterLogger(logger log.Logger) RateLimiterConfigOpt {
	return func(config *RateLimiterConfig) {
		config.Logger = logger
	}
}

// WithCommandsPerMinute sets the number of commands per minute that the Gateway will allow.
func WithCommandsPerMinute(commandsPerMinute int) RateLimiterConfigOpt {
	return func(config *RateLimiterConfig) {
		config.CommandsPerMinute = commandsPerMinute
	}
}
