package grate

import (
	"github.com/disgoorg/log"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Logger:            log.Default(),
		CommandsPerMinute: 120,
	}
}

// Config lets you configure your Gateway instance.
type Config struct {
	Logger            log.Logger
	CommandsPerMinute int
}

// ConfigOpt is a type alias for a function that takes a Config and is used to configure your Server.
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger sets the Logger for the Gateway.
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithCommandsPerMinute sets the number of commands per minute that the Gateway will allow.
func WithCommandsPerMinute(commandsPerMinute int) ConfigOpt {
	return func(config *Config) {
		config.CommandsPerMinute = commandsPerMinute
	}
}
