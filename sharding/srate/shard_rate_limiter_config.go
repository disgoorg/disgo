package srate

import (
	"github.com/disgoorg/log"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Logger:         log.Default(),
		MaxConcurrency: 1,
	}
}

// Config lets you configure your Limiter instance.
type Config struct {
	Logger         log.Logger
	MaxConcurrency int
}

// ConfigOpt is a type alias for a function that takes a Config and is used to configure your Server.
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger sets the logger for the Limiter.
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithMaxConcurrency sets the maximum number of concurrent identifies in 5 seconds.
func WithMaxConcurrency(maxConcurrency int) ConfigOpt {
	return func(config *Config) {
		config.MaxConcurrency = maxConcurrency
	}
}
