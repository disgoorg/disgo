package gateway

import (
	"log/slog"
	"time"
)

func defaultIdentifyRateLimiterConfig() identifyRateLimiterConfig {
	return identifyRateLimiterConfig{
		Logger:         slog.Default(),
		MaxConcurrency: DefaultMaxConcurrency,
		Wait:           5 * time.Second,
	}
}

type identifyRateLimiterConfig struct {
	Logger         *slog.Logger
	MaxConcurrency int
	Wait           time.Duration
}

// IdentifyRateLimiterConfigOpt is a type alias for a function that takes a identifyRateLimiterConfig and is used to configure your Server.
type IdentifyRateLimiterConfigOpt func(config *identifyRateLimiterConfig)

func (c *identifyRateLimiterConfig) apply(opts []IdentifyRateLimiterConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "sharding_rate_limiter"))
}

// WithIdentifyRateLimiterLogger sets the logger for the RateLimiter.
func WithIdentifyRateLimiterLogger(logger *slog.Logger) IdentifyRateLimiterConfigOpt {
	return func(config *identifyRateLimiterConfig) {
		config.Logger = logger
	}
}

// WithIdentifyMaxConcurrency sets the maximum number of concurrent identifies allowed per window duration (default is 5 seconds, configurable via WithIdentifyWait).
func WithIdentifyMaxConcurrency(maxConcurrency int) IdentifyRateLimiterConfigOpt {
	return func(config *identifyRateLimiterConfig) {
		config.MaxConcurrency = maxConcurrency
	}
}

// WithIdentifyWait sets the duration to wait in between identifying shards.
func WithIdentifyWait(identifyWait time.Duration) IdentifyRateLimiterConfigOpt {
	return func(config *identifyRateLimiterConfig) {
		config.Wait = identifyWait
	}
}
