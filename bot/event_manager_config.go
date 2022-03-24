package bot

import "github.com/disgoorg/disgo/discord"

func DefaultEventManagerConfig() *EventManagerConfig {
	return &EventManagerConfig{}
}

type EventManagerConfig struct {
	EventListeners     []EventListener
	RawEventsEnabled   bool
	AsyncEventsEnabled bool

	GatewayHandlers   map[discord.GatewayEventType]GatewayEventHandler
	HTTPServerHandler HTTPServerEventHandler
}

type EventManagerConfigOpt func(config *EventManagerConfig)

func (c *EventManagerConfig) Apply(opts []EventManagerConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithListeners(listeners ...EventListener) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.EventListeners = append(config.EventListeners, listeners...)
	}
}

func WithRawEventsEnabled() EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.RawEventsEnabled = true
	}
}

func WithAsyncEventsEnabled() EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.AsyncEventsEnabled = true
	}
}

func WithGatewayHandlers(handlers map[discord.GatewayEventType]GatewayEventHandler) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.GatewayHandlers = handlers
	}
}

func WithHTTPServerHandler(handler HTTPServerEventHandler) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.HTTPServerHandler = handler
	}
}
