package voice

import (
	"github.com/disgoorg/log"
)

func DefaultConnectionConfig() *ConnectionConfig {
	return &ConnectionConfig{
		Logger:            log.Default(),
		GatewayCreateFunc: NewGateway,
		UDPConnCreateFunc: NewUDP,
	}
}

type ConnectionConfig struct {
	Logger log.Logger

	GatewayCreateFunc GatewayCreateFunc
	GatewayConfigOpts []GatewayConfigOpt

	UDPConnCreateFunc UDPCreateFunc
	UDPConnConfigOpts []UDPConfigOpt

	EventHandlerFunc EventHandlerFunc
}

type ConnectionConfigOpt func(ConnectionConfig *ConnectionConfig)

func (c *ConnectionConfig) Apply(opts []ConnectionConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithConnectionLogger(logger log.Logger) ConnectionConfigOpt {
	return func(ConnectionConfig *ConnectionConfig) {
		ConnectionConfig.Logger = logger
	}
}

func WithConnectionEventHandlerFunc(eventHandlerFunc EventHandlerFunc) ConnectionConfigOpt {
	return func(ConnectionConfig *ConnectionConfig) {
		ConnectionConfig.EventHandlerFunc = eventHandlerFunc
	}
}
