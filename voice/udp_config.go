package voice

import (
	"net"
	"time"

	"github.com/disgoorg/log"
)

func DefaultUDPConnConfig() *UDPConnConfig {
	return &UDPConnConfig{
		Logger:            log.Default(),
		UDPConnCreateFunc: NewUDPConn,
		Dialer: &net.Dialer{
			Timeout: 30 * time.Second,
		},
	}
}

type UDPConnConfig struct {
	Logger            log.Logger
	UDPConnCreateFunc UDPConnCreateFunc
	Dialer            *net.Dialer
}

type UDPConnConfigOpt func(config *UDPConnConfig)

func (c *UDPConnConfig) Apply(opts []UDPConnConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithUDPConnLogger(logger log.Logger) UDPConnConfigOpt {
	return func(config *UDPConnConfig) {
		config.Logger = logger
	}
}
