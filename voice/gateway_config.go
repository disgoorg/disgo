package voice

import (
	"github.com/disgoorg/log"
	"github.com/gorilla/websocket"
)

func DefaultGatewayConfig() *GatewayConfig {
	return &GatewayConfig{
		Logger:            log.Default(),
		Dialer:            websocket.DefaultDialer,
		AutoReconnect:     true,
		MaxReconnectTries: 10,
	}
}

type GatewayConfig struct {
	Logger            log.Logger
	Dialer            *websocket.Dialer
	AutoReconnect     bool
	MaxReconnectTries int
}

type GatewayConfigOpt func(config *GatewayConfig)

func (c *GatewayConfig) Apply(opts []GatewayConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithGatewayLogger(logger log.Logger) GatewayConfigOpt {
	return func(config *GatewayConfig) {
		config.Logger = logger
	}
}

func WithGatewayDialer(dialer *websocket.Dialer) GatewayConfigOpt {
	return func(config *GatewayConfig) {
		config.Dialer = dialer
	}
}

func WithGatewayAutoReconnect(autoReconnect bool) GatewayConfigOpt {
	return func(config *GatewayConfig) {
		config.AutoReconnect = autoReconnect
	}
}
