package srate

import (
	"github.com/DisgoOrg/log"
)

func DefaultConfig() *Config {
	return &Config{
		Logger:         log.Default(),
		MaxConcurrency: 1,
		StartupDelay:   5,
	}
}

type Config struct {
	Logger         log.Logger
	MaxConcurrency int
	StartupDelay   int
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithMaxConcurrency(maxConcurrency int) ConfigOpt {
	return func(config *Config) {
		config.MaxConcurrency = maxConcurrency
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithStartupDelay(startupDelay int) ConfigOpt {
	return func(config *Config) {
		config.StartupDelay = startupDelay
	}
}
