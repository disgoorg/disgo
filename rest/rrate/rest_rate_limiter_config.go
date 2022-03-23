package rrate

import (
	"github.com/disgoorg/log"
)

// DefaultConfig is the configuration which is used by default
func DefaultConfig() *Config {
	return &Config{
		Logger:     log.Default(),
		MaxRetries: 10,
	}
}

// Config is the configuration for the rate limiter
type Config struct {
	Logger     log.Logger
	MaxRetries int
}

// ConfigOpt can be used to supply optional parameters to NewLimiter
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger applies a custom logger to the rest rate limiter
//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithMaxRetries tells the rest rate limiter to retry the request up to the specified number of times if it encounters a 429 response
//goland:noinspection GoUnusedExportedFunction
func WithMaxRetries(maxRetries int) ConfigOpt {
	return func(config *Config) {
		config.MaxRetries = maxRetries
	}
}
