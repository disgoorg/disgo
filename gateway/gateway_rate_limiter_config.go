package gateway

import (
	"log/slog"
)

// DefaultRateLimiterConfig returns a RateLimiterConfig with sensible defaults.
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		Logger:            slog.Default(),
		CommandsPerMinute: CommandsPerMinute,
	}
}

// RateLimiterConfig lets you configure your Gateway instance.
type RateLimiterConfig struct {
	Logger            *slog.Logger
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
func WithRateLimiterLogger(logger *slog.Logger) RateLimiterConfigOpt {
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
