package rate

import (
	"github.com/DisgoOrg/log"
)

var DefaultConfig = Config{
	Logger:         log.Default(),
	MaxConcurrency: 1,
}

type Config struct {
	Logger         log.Logger
	MaxConcurrency int
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

func WithMaxConcurrency(maxConcurrency int) ConfigOpt {
	return func(config *Config) {
		config.MaxConcurrency = maxConcurrency
	}
}
