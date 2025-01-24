package voice

import (
	"log/slog"
)

// DefaultConnConfig returns a ConnConfig with sensible defaults.
func DefaultConnConfig() *ConnConfig {
	return &ConnConfig{
		Logger:                  slog.Default(),
		GatewayCreateFunc:       NewGateway,
		UDPConnCreateFunc:       NewUDPConn,
		AudioSenderCreateFunc:   NewAudioSender,
		AudioReceiverCreateFunc: NewAudioReceiver,
	}
}

// ConnConfig is used to configure a Conn.
type ConnConfig struct {
	Logger *slog.Logger

	GatewayCreateFunc GatewayCreateFunc
	GatewayConfigOpts []GatewayConfigOpt

	UDPConnCreateFunc UDPConnCreateFunc
	UDPConnConfigOpts []UDPConnConfigOpt

	AudioSenderCreateFunc   AudioSenderCreateFunc
	AudioReceiverCreateFunc AudioReceiverCreateFunc

	EventHandlerFunc EventHandlerFunc
}

// ConnConfigOpt is used to functionally configure a ConnConfig.
type ConnConfigOpt func(connConfig *ConnConfig)

// Apply applies the ConnConfigOpt(s) to the ConnConfig.
func (c *ConnConfig) Apply(opts []ConnConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithConnLogger sets the Conn(s) used Logger.
func WithConnLogger(logger *slog.Logger) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.Logger = logger
	}
}

// WithConnGatewayCreateFunc sets the Conn(s) used GatewayCreateFunc.
func WithConnGatewayCreateFunc(gatewayCreateFunc GatewayCreateFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}

// WithConnGatewayConfigOpts sets the Conn(s) used GatewayConfigOpt(s).
func WithConnGatewayConfigOpts(opts ...GatewayConfigOpt) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

// WithUDPConnCreateFunc sets the Conn(s) used UDPConnCreateFunc.
func WithUDPConnCreateFunc(udpConnCreateFunc UDPConnCreateFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.UDPConnCreateFunc = udpConnCreateFunc
	}
}

// WithUDPConnConfigOpts sets the Conn(s) used UDPConnConfigOpt(s).
func WithUDPConnConfigOpts(opts ...UDPConnConfigOpt) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.UDPConnConfigOpts = append(config.UDPConnConfigOpts, opts...)
	}
}

// WithConnAudioSenderCreateFunc sets the Conn(s) used AudioSenderCreateFunc.
func WithConnAudioSenderCreateFunc(audioSenderCreateFunc AudioSenderCreateFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.AudioSenderCreateFunc = audioSenderCreateFunc
	}
}

// WithConnAudioReceiverCreateFunc sets the Conn(s) used AudioReceiverCreateFunc.
func WithConnAudioReceiverCreateFunc(audioReceiverCreateFunc AudioReceiverCreateFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.AudioReceiverCreateFunc = audioReceiverCreateFunc
	}
}

// WithConnEventHandlerFunc sets the Conn(s) used EventHandlerFunc.
func WithConnEventHandlerFunc(eventHandlerFunc EventHandlerFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.EventHandlerFunc = eventHandlerFunc
	}
}
