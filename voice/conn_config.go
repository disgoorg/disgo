package voice

import (
	"time"

	"github.com/disgoorg/disgo/voice/gateway"
	"github.com/disgoorg/disgo/voice/udp"
	"github.com/disgoorg/log"
)

func DefaultConnConfig() *ConnConfig {
	return &ConnConfig{
		Logger:                  log.Default(),
		GatewayCreateFunc:       gateway.New,
		UDPConnCreateFunc:       udp.NewConn,
		AudioSenderCreateFunc:   NewAudioSender,
		AudioReceiverCreateFunc: NewAudioReceiver,
		ReconnectTimeout:        5 * time.Second,
	}
}

type ConnConfig struct {
	Logger log.Logger

	GatewayCreateFunc gateway.CreateFunc
	GatewayConfigOpts []gateway.ConfigOpt

	UDPConnCreateFunc udp.ConnCreateFunc
	UDPConnConfigOpts []udp.ConnConfigOpt

	AudioSenderCreateFunc   AudioSenderCreateFunc
	AudioReceiverCreateFunc AudioReceiverCreateFunc

	EventHandlerFunc gateway.EventHandlerFunc
	ReconnectTimeout time.Duration
}

type ConnConfigOpt func(connConfig *ConnConfig)

func (c *ConnConfig) Apply(opts []ConnConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithConnLogger(logger log.Logger) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.Logger = logger
	}
}

func WithConnGatewayCreateFunc(gatewayCreateFunc gateway.CreateFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}

func WithConnGatewayConfigOpts(opts ...gateway.ConfigOpt) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

func WithUDPConnCreateFunc(udpConnCreateFunc udp.ConnCreateFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.UDPConnCreateFunc = udpConnCreateFunc
	}
}

func WithUDPConnConfigOpts(opts ...udp.ConnConfigOpt) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.UDPConnConfigOpts = append(config.UDPConnConfigOpts, opts...)
	}
}

func WithConnAudioSenderCreateFunc(audioSenderCreateFunc AudioSenderCreateFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.AudioSenderCreateFunc = audioSenderCreateFunc
	}
}

func WithConnAudioReceiverCreateFunc(audioReceiverCreateFunc AudioReceiverCreateFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.AudioReceiverCreateFunc = audioReceiverCreateFunc
	}
}

func WithConnEventHandlerFunc(eventHandlerFunc gateway.EventHandlerFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.EventHandlerFunc = eventHandlerFunc
	}
}
