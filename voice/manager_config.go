package voice

import (
	"log/slog"

	"github.com/disgoorg/godave"
)

func defaultManagerConfig() managerConfig {
	return managerConfig{
		Logger:         slog.Default(),
		ConnCreateFunc: NewConn,
	}
}

type managerConfig struct {
	Logger *slog.Logger

	ConnCreateFunc ConnCreateFunc
	ConnOpts       []ConnConfigOpt
}

// ManagerConfigOpt is used to functionally configure a managerConfig.
type ManagerConfigOpt func(config *managerConfig)

func (c *managerConfig) apply(opts []ManagerConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "voice"))
}

// WithLogger sets the logger for the webhook client.
func WithLogger(logger *slog.Logger) ManagerConfigOpt {
	return func(config *managerConfig) {
		config.Logger = logger
	}
}

// WithConnCreateFunc sets the ConnCreateFunc for the Manager.
func WithConnCreateFunc(connectionCreateFunc ConnCreateFunc) ManagerConfigOpt {
	return func(config *managerConfig) {
		config.ConnCreateFunc = connectionCreateFunc
	}
}

// WithConnConfigOpts sets the ConnConfigOpt(s) for the Manager.
func WithConnConfigOpts(opts ...ConnConfigOpt) ManagerConfigOpt {
	return func(config *managerConfig) {
		config.ConnOpts = append(config.ConnOpts, opts...)
	}
}

// WithDaveSessionCreateFunc sets the godave.SessionCreateFunc for the Manager.
func WithDaveSessionCreateFunc(sessionCreateFunc godave.SessionCreateFunc) ManagerConfigOpt {
	return WithConnConfigOpts(WithConnDaveSessionCreateFunc(sessionCreateFunc))
}
