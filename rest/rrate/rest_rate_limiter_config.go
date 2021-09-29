package rrate

import (
	"github.com/DisgoOrg/log"
)

var DefaultConfig = Config{
	MaxRetries: 10,
}

type Config struct {
	Logger     log.Logger
	MaxRetries int
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

func WithMaxRetries(maxRetries int) ConfigOpt {
	return func(config *Config) {
		config.MaxRetries = maxRetries
	}
}
