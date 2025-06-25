package rest

import (
	"log/slog"
	"time"
)

func defaultRateLimiterConfig() rateLimiterConfig {
	return rateLimiterConfig{
		Logger:          slog.Default(),
		MaxRetries:      MaxRetries,
		CleanupInterval: CleanupInterval,
	}
}

type rateLimiterConfig struct {
	Logger          *slog.Logger
	MaxRetries      int
	CleanupInterval time.Duration
}

// RateLimiterConfigOpt can be used to supply optional parameters to NewRateLimiter.
type RateLimiterConfigOpt func(config *rateLimiterConfig)

func (c *rateLimiterConfig) apply(opts []RateLimiterConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "rest_rate_limiter"))
}

// WithRateLimiterLogger applies a custom logger to the rest rate limiter.
func WithRateLimiterLogger(logger *slog.Logger) RateLimiterConfigOpt {
	return func(config *rateLimiterConfig) {
		config.Logger = logger
	}
}

// WithMaxRetries tells the rest rate limiter to retry the request up to the specified number of times if it encounters a 429 response.
func WithMaxRetries(maxRetries int) RateLimiterConfigOpt {
	return func(config *rateLimiterConfig) {
		config.MaxRetries = maxRetries
	}
}

// WithCleanupInterval tells the rest rate limiter how often to clean up the rate limiter buckets.
func WithCleanupInterval(cleanupInterval time.Duration) RateLimiterConfigOpt {
	return func(config *rateLimiterConfig) {
		config.CleanupInterval = cleanupInterval
	}
}
