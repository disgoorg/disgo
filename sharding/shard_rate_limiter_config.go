package sharding

import (
	"log/slog"
	"time"
)

func defaultRateLimiterConfig() rateLimiterConfig {
	return rateLimiterConfig{
		Logger:         slog.Default(),
		MaxConcurrency: MaxConcurrency,
		IdentifyWait:   5 * time.Second,
	}
}

type rateLimiterConfig struct {
	Logger         *slog.Logger
	MaxConcurrency int
	IdentifyWait   time.Duration
}

// RateLimiterConfigOpt is a type alias for a function that takes a rateLimiterConfig and is used to configure your Server.
type RateLimiterConfigOpt func(config *rateLimiterConfig)

func (c *rateLimiterConfig) apply(opts []RateLimiterConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "sharding_rate_limiter"))
}

// WithRateLimiterLogger sets the logger for the RateLimiter.
func WithRateLimiterLogger(logger *slog.Logger) RateLimiterConfigOpt {
	return func(config *rateLimiterConfig) {
		config.Logger = logger
	}
}

// WithMaxConcurrency sets the maximum number of concurrent identifies in 5 seconds.
func WithMaxConcurrency(maxConcurrency int) RateLimiterConfigOpt {
	return func(config *rateLimiterConfig) {
		config.MaxConcurrency = maxConcurrency
	}
}

// WithIdentifyWait sets the duration to wait in between identifying shards.
func WithIdentifyWait(identifyWait time.Duration) RateLimiterConfigOpt {
	return func(config *rateLimiterConfig) {
		config.IdentifyWait = identifyWait
	}
}
