package voice

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

// DefaultGatewayConfig returns a GatewayConfig with sensible defaults.
func DefaultGatewayConfig() *GatewayConfig {
	return &GatewayConfig{
		Logger:        slog.Default(),
		Dialer:        websocket.DefaultDialer,
		AutoReconnect: true,
	}
}

// GatewayConfig is used to configure a Gateway.
type GatewayConfig struct {
	Logger        *slog.Logger
	Dialer        *websocket.Dialer
	AutoReconnect bool
}

// GatewayConfigOpt is used to functionally configure a GatewayConfig.
type GatewayConfigOpt func(config *GatewayConfig)

// Apply applies the GatewayConfigOpt(s) to the GatewayConfig.
func (c *GatewayConfig) Apply(opts []GatewayConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithGatewayLogger sets the Gateway(s) used Logger.
func WithGatewayLogger(logger *slog.Logger) GatewayConfigOpt {
	return func(config *GatewayConfig) {
		config.Logger = logger
	}
}

// WithGatewayDialer sets the Gateway(s) used websocket.Dialer.
func WithGatewayDialer(dialer *websocket.Dialer) GatewayConfigOpt {
	return func(config *GatewayConfig) {
		config.Dialer = dialer
	}
}

// WithGatewayAutoReconnect sets the Gateway(s) used AutoReconnect.
func WithGatewayAutoReconnect(autoReconnect bool) GatewayConfigOpt {
	return func(config *GatewayConfig) {
		config.AutoReconnect = autoReconnect
	}
}
