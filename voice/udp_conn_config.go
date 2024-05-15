package voice

import (
	"log/slog"
	"net"
)

func DefaultUDPConnConfig() UDPConnConfig {
	return UDPConnConfig{
		Logger: slog.Default(),
		Dialer: &net.Dialer{
			Timeout: UDPTimeout,
		},
	}
}

type UDPConnConfig struct {
	Logger *slog.Logger
	Dialer *net.Dialer
}

type UDPConnConfigOpt func(config *UDPConnConfig)

func (c *UDPConnConfig) Apply(opts []UDPConnConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithUDPConnLogger(logger *slog.Logger) UDPConnConfigOpt {
	return func(config *UDPConnConfig) {
		config.Logger = logger
	}
}

func WithUDPConnDialer(dialer *net.Dialer) UDPConnConfigOpt {
	return func(config *UDPConnConfig) {
		config.Dialer = dialer
	}
}
