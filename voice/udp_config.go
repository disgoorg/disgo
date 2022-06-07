package voice

import (
	"net"
	"time"

	"github.com/disgoorg/log"
)

func DefaultUDPConfig() *UDPConfig {
	return &UDPConfig{
		Logger:        log.Default(),
		UDPCreateFunc: NewUDP,
		Dialer: &net.Dialer{
			Timeout: 30 * time.Second,
		},
	}
}

type UDPConfig struct {
	Logger        log.Logger
	UDPCreateFunc UDPCreateFunc
	Dialer        *net.Dialer
}

type UDPConfigOpt func(config *UDPConfig)

func (c *UDPConfig) Apply(opts []UDPConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithUDPLogger(logger log.Logger) UDPConfigOpt {
	return func(config *UDPConfig) {
		config.Logger = logger
	}
}
