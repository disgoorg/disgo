package voice

import (
	"log/slog"
	"net"
)

func defaultUDPConnConfig() udpConnConfig {
	return udpConnConfig{
		Logger: slog.Default(),
		Dialer: &net.Dialer{
			Timeout: UDPTimeout,
		},
	}
}

type udpConnConfig struct {
	Logger *slog.Logger
	Dialer *net.Dialer
}

// UDPConnConfigOpt is a function that modifies the udpConnConfig.
type UDPConnConfigOpt func(config *udpConnConfig)

func (c *udpConnConfig) apply(opts []UDPConnConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "voice_conn_udp_conn"))
}

func WithUDPConnLogger(logger *slog.Logger) UDPConnConfigOpt {
	return func(config *udpConnConfig) {
		config.Logger = logger
	}
}

func WithUDPConnDialer(dialer *net.Dialer) UDPConnConfigOpt {
	return func(config *udpConnConfig) {
		config.Dialer = dialer
	}
}
