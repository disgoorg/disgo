package sharding

import (
	"github.com/disgoorg/log"
)

// DefaultRateLimiterConfig returns a RateLimiterConfig with sensible defaults.
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		Logger:         log.Default(),
		MaxConcurrency: 1,
	}
}

// RateLimiterConfig lets you configure your RateLimiter instance.
type RateLimiterConfig struct {
	Logger         log.Logger
	MaxConcurrency int
}

// RateLimiterConfigOpt is a type alias for a function that takes a RateLimiterConfig and is used to configure your Server.
type RateLimiterConfigOpt func(config *RateLimiterConfig)

// Apply applies the given RateLimiterConfigOpt(s) to the RateLimiterConfig
func (c *RateLimiterConfig) Apply(opts []RateLimiterConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithRateLimiterLogger sets the logger for the RateLimiter.
func WithRateLimiterLogger(logger log.Logger) RateLimiterConfigOpt {
	return func(config *RateLimiterConfig) {
		config.Logger = logger
	}
}

// WithMaxConcurrency sets the maximum number of concurrent identifies in 5 seconds.
func WithMaxConcurrency(maxConcurrency int) RateLimiterConfigOpt {
	return func(config *RateLimiterConfig) {
		config.MaxConcurrency = maxConcurrency
	}
}
