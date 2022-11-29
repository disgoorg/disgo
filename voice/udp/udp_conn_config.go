package udp

import (
	"net"
	"time"

	"github.com/disgoorg/log"
)

func DefaultConfig() *ConnConfig {
	return &ConnConfig{
		Logger: log.Default(),
		Dialer: &net.Dialer{
			Timeout: 30 * time.Second,
		},
	}
}

type ConnConfig struct {
	Logger log.Logger
	Dialer *net.Dialer
}

type ConnConfigOpt func(config *ConnConfig)

func (c *ConnConfig) Apply(opts []ConnConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithLogger(logger log.Logger) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.Logger = logger
	}
}

func WithDialer(dialer *net.Dialer) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.Dialer = dialer
	}
}
