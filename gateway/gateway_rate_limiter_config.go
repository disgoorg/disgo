package gateway

import (
	"log/slog"
)

func defaultRateLimiterConfig() rateLimiterConfig {
	return rateLimiterConfig{
		Logger:               slog.Default(),
		CommandsPerMinute:    CommandsPerMinute,
		ReservedCommandSlots: ReservedCommandSlots,
	}
}

type rateLimiterConfig struct {
	Logger               *slog.Logger
	CommandsPerMinute    int
	ReservedCommandSlots int
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

// WithReservedCommandSlots sets the number of reserved slots for Gateway priority events, like heartbeats.
func WithReservedCommandSlots(reservedCommandSlots int) RateLimiterConfigOpt {
	return func(config *rateLimiterConfig) {
		config.ReservedCommandSlots = reservedCommandSlots
	}
}
