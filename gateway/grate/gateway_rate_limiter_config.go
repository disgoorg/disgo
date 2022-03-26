package grate

import (
	"github.com/disgoorg/log"
)

func DefaultConfig() *Config {
	return &Config{
		Logger:            log.Default(),
		CommandsPerMinute: 120,
	}
}

type Config struct {
	Logger            log.Logger
	CommandsPerMinute int
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

func WithCommandsPerMinute(commandsPerMinute int) ConfigOpt {
	return func(config *Config) {
		config.CommandsPerMinute = commandsPerMinute
	}
}
