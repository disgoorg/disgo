package voice

import (
	"github.com/disgoorg/log"
)

func DefaultConfig() *Config {
	return &Config{
		Logger:            log.Default(),
		GatewayCreateFunc: NewGateway,
		UDPConnCreateFunc: NewUDPConn,
	}
}

type Config struct {
	Logger log.Logger

	GatewayCreateFunc GatewayCreateFunc
	GatewayConfigOpts []GatewayConfigOpt

	UDPConnCreateFunc UDPConnCreateFunc
	UDPConnConfigOpts []UDPConnConfigOpt
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger sets the logger for the webhook client
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}
