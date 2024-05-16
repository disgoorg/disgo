package rest

import (
	"log/slog"
	"time"
)

// DefaultRateLimiterConfig is the configuration which is used by default.
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		Logger:          slog.Default(),
		MaxRetries:      MaxRetries,
		CleanupInterval: CleanupInterval,
	}
}

// RateLimiterConfig is the configuration for the rate limiter.
type RateLimiterConfig struct {
	Logger          *slog.Logger
	MaxRetries      int
	CleanupInterval time.Duration
}

// RateLimiterConfigOpt can be used to supply optional parameters to NewRateLimiter.
type RateLimiterConfigOpt func(config *RateLimiterConfig)

// Apply applies the given RateLimiterConfigOpt(s) to the RateLimiterConfig.
func (c *RateLimiterConfig) Apply(opts []RateLimiterConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithRateLimiterLogger applies a custom logger to the rest rate limiter.
func WithRateLimiterLogger(logger *slog.Logger) RateLimiterConfigOpt {
	return func(config *RateLimiterConfig) {
		config.Logger = logger
	}
}

// WithMaxRetries tells the rest rate limiter to retry the request up to the specified number of times if it encounters a 429 response.
func WithMaxRetries(maxRetries int) RateLimiterConfigOpt {
	return func(config *RateLimiterConfig) {
		config.MaxRetries = maxRetries
	}
}

// WithCleanupInterval tells the rest rate limiter how often to clean up the rate limiter buckets.
func WithCleanupInterval(cleanupInterval time.Duration) RateLimiterConfigOpt {
	return func(config *RateLimiterConfig) {
		config.CleanupInterval = cleanupInterval
	}
}
