package voice

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

func defaultGatewayConfig() gatewayConfig {
	return gatewayConfig{
		Logger:        slog.Default(),
		Dialer:        websocket.DefaultDialer,
		AutoReconnect: true,
	}
}

type gatewayConfig struct {
	Logger        *slog.Logger
	Dialer        *websocket.Dialer
	AutoReconnect bool
}

// GatewayConfigOpt is used to functionally configure a gatewayConfig.
type GatewayConfigOpt func(config *gatewayConfig)

func (c *gatewayConfig) apply(opts []GatewayConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "voice_conn_gateway"))
}

// WithGatewayLogger sets the Gateway(s) used Logger.
func WithGatewayLogger(logger *slog.Logger) GatewayConfigOpt {
	return func(config *gatewayConfig) {
		config.Logger = logger
	}
}

// WithGatewayDialer sets the Gateway(s) used websocket.Dialer.
func WithGatewayDialer(dialer *websocket.Dialer) GatewayConfigOpt {
	return func(config *gatewayConfig) {
		config.Dialer = dialer
	}
}

// WithGatewayAutoReconnect sets the Gateway(s) used AutoReconnect.
func WithGatewayAutoReconnect(autoReconnect bool) GatewayConfigOpt {
	return func(config *gatewayConfig) {
		config.AutoReconnect = autoReconnect
	}
}
