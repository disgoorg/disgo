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

func WithConnectionGatewayCreateFunc(gatewayCreateFunc GatewayCreateFunc) ConnectionConfigOpt {
	return func(ConnectionConfig *ConnectionConfig) {
		ConnectionConfig.GatewayCreateFunc = gatewayCreateFunc
	}
}

func WithConnectionGatewayConfigOpts(opts ...GatewayConfigOpt) ConnectionConfigOpt {
	return func(ConnectionConfig *ConnectionConfig) {
		ConnectionConfig.GatewayConfigOpts = append(ConnectionConfig.GatewayConfigOpts, opts...)
	}
}

func WithConnectionUDPCreateFunc(udpCreateFunc UDPCreateFunc) ConnectionConfigOpt {
	return func(ConnectionConfig *ConnectionConfig) {
		ConnectionConfig.UDPConnCreateFunc = udpCreateFunc
	}
}

func WithConnectionUDPConfigOpts(opts ...UDPConfigOpt) ConnectionConfigOpt {
	return func(ConnectionConfig *ConnectionConfig) {
		ConnectionConfig.UDPConnConfigOpts = append(ConnectionConfig.UDPConnConfigOpts, opts...)
	}
}

func WithConnectionEventHandlerFunc(eventHandlerFunc EventHandlerFunc) ConnectionConfigOpt {
	return func(ConnectionConfig *ConnectionConfig) {
		ConnectionConfig.EventHandlerFunc = eventHandlerFunc
	}
}
