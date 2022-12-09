package voice

import (
	"github.com/disgoorg/disgo/voice/voicegateway"
	"github.com/disgoorg/disgo/voice/voiceudp"
	"github.com/disgoorg/log"
)

func DefaultConnConfig() *ConnConfig {
	return &ConnConfig{
		Logger:                  log.Default(),
		GatewayCreateFunc:       voicegateway.New,
		UDPConnCreateFunc:       voiceudp.NewConn,
		AudioSenderCreateFunc:   NewAudioSender,
		AudioReceiverCreateFunc: NewAudioReceiver,
	}
}

type ConnConfig struct {
	Logger log.Logger

	GatewayCreateFunc voicegateway.CreateFunc
	GatewayConfigOpts []voicegateway.ConfigOpt

	UDPConnCreateFunc voiceudp.ConnCreateFunc
	UDPConnConfigOpts []voiceudp.ConnConfigOpt

	AudioSenderCreateFunc   AudioSenderCreateFunc
	AudioReceiverCreateFunc AudioReceiverCreateFunc

	EventHandlerFunc voicegateway.EventHandlerFunc
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

func WithConnGatewayCreateFunc(gatewayCreateFunc voicegateway.CreateFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.GatewayCreateFunc = gatewayCreateFunc
	}
}

func WithConnGatewayConfigOpts(opts ...voicegateway.ConfigOpt) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

func WithUDPConnCreateFunc(udpConnCreateFunc voiceudp.ConnCreateFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.UDPConnCreateFunc = udpConnCreateFunc
	}
}

func WithUDPConnConfigOpts(opts ...voiceudp.ConnConfigOpt) ConnConfigOpt {
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

func WithConnEventHandlerFunc(eventHandlerFunc voicegateway.EventHandlerFunc) ConnConfigOpt {
	return func(config *ConnConfig) {
		config.EventHandlerFunc = eventHandlerFunc
	}
}
