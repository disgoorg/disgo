package gateway

import (
	"log/slog"
)

func defaultRateLimiterConfig() rateLimiterConfig {
	return rateLimiterConfig{
		Logger:            slog.Default(),
		CommandsPerMinute: CommandsPerMinute,
	}
}

type rateLimiterConfig struct {
	Logger            *slog.Logger
	CommandsPerMinute int
}

// RateLimiterConfigOpt is a type alias for a function that takes a rateLimiterConfig and is used to configure your Server.
type RateLimiterConfigOpt func(config *rateLimiterConfig)

func (c *rateLimiterConfig) apply(opts []RateLimiterConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "gateway_rate_limiter"))
}

// WithRateLimiterLogger sets the Logger for the Gateway.
func WithRateLimiterLogger(logger *slog.Logger) RateLimiterConfigOpt {
	return func(config *rateLimiterConfig) {
		config.Logger = logger
	}
}

// WithCommandsPerMinute sets the number of commands per minute that the Gateway will allow.
func WithCommandsPerMinute(commandsPerMinute int) RateLimiterConfigOpt {
	return func(config *rateLimiterConfig) {
		config.CommandsPerMinute = commandsPerMinute
	}
}
