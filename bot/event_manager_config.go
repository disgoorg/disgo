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

//goland:noinspection GoUnusedExportedFunction
func WithEventListeners(listeners ...EventListener) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.EventListeners = append(config.EventListeners, listeners...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRawEventsEnabled() EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.RawEventsEnabled = true
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithAsyncEventsEnabled() EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.AsyncEventsEnabled = true
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayHandlers(handlers map[discord.GatewayEventType]GatewayEventHandler) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.GatewayHandlers = handlers
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithHTTPServerHandler(handler HTTPServerEventHandler) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.HTTPServerHandler = handler
	}
}
