package voice

import (
	"log/slog"
)

func defaultConnConfig() connConfig {
	return connConfig{
		Logger:                  slog.Default(),
		GatewayCreateFunc:       NewGateway,
		UDPConnCreateFunc:       NewUDPConn,
		AudioSenderCreateFunc:   NewAudioSender,
		AudioReceiverCreateFunc: NewAudioReceiver,
	}
}

type connConfig struct {
	Logger *slog.Logger

	GatewayCreateFunc GatewayCreateFunc
	GatewayConfigOpts []GatewayConfigOpt

	UDPConnCreateFunc UDPConnCreateFunc
	UDPConnConfigOpts []UDPConnConfigOpt

	AudioSenderCreateFunc   AudioSenderCreateFunc
	AudioReceiverCreateFunc AudioReceiverCreateFunc

	EventHandlerFunc EventHandlerFunc
}

// ConnConfigOpt is used to functionally configure a connConfig.
type ConnConfigOpt func(config *connConfig)

func (c *connConfig) apply(opts []ConnConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "voice_conn"))
}

// WithConnLogger sets the Conn(s) used Logger.
func WithConnLogger(logger *slog.Logger) ConnConfigOpt {
	return func(config *connConfig) {
		config.Logger = logger
	}
}

// WithConnGatewayCreateFunc sets the Conn(s) used GatewayCreateFunc.
func WithConnGatewayCreateFunc(gatewayCreateFunc GatewayCreateFunc) ConnConfigOpt {
	return func(config *connConfig) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}

// WithConnGatewayConfigOpts sets the Conn(s) used GatewayConfigOpt(s).
func WithConnGatewayConfigOpts(opts ...GatewayConfigOpt) ConnConfigOpt {
	return func(config *connConfig) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

// WithUDPConnCreateFunc sets the Conn(s) used UDPConnCreateFunc.
func WithUDPConnCreateFunc(udpConnCreateFunc UDPConnCreateFunc) ConnConfigOpt {
	return func(config *connConfig) {
		config.UDPConnCreateFunc = udpConnCreateFunc
	}
}

// WithUDPConnConfigOpts sets the Conn(s) used UDPConnConfigOpt(s).
func WithUDPConnConfigOpts(opts ...UDPConnConfigOpt) ConnConfigOpt {
	return func(config *connConfig) {
		config.UDPConnConfigOpts = append(config.UDPConnConfigOpts, opts...)
	}
}

// WithConnAudioSenderCreateFunc sets the Conn(s) used AudioSenderCreateFunc.
func WithConnAudioSenderCreateFunc(audioSenderCreateFunc AudioSenderCreateFunc) ConnConfigOpt {
	return func(config *connConfig) {
		config.AudioSenderCreateFunc = audioSenderCreateFunc
	}
}

// WithConnAudioReceiverCreateFunc sets the Conn(s) used AudioReceiverCreateFunc.
func WithConnAudioReceiverCreateFunc(audioReceiverCreateFunc AudioReceiverCreateFunc) ConnConfigOpt {
	return func(config *connConfig) {
		config.AudioReceiverCreateFunc = audioReceiverCreateFunc
	}
}

// WithConnEventHandlerFunc sets the Conn(s) used EventHandlerFunc.
func WithConnEventHandlerFunc(eventHandlerFunc EventHandlerFunc) ConnConfigOpt {
	return func(config *connConfig) {
		config.EventHandlerFunc = eventHandlerFunc
	}
}
