package bot

import "github.com/disgoorg/disgo/discord"

// DefaultEventManagerConfig returns a new EventManagerConfig with all default values.
func DefaultEventManagerConfig() *EventManagerConfig {
	return &EventManagerConfig{}
}

// EventManagerConfig can be used to configure the EventManager.
type EventManagerConfig struct {
	EventListeners     []EventListener
	RawEventsEnabled   bool
	AsyncEventsEnabled bool

	GatewayHandlers   map[discord.GatewayEventType]GatewayEventHandler
	HTTPServerHandler HTTPServerEventHandler
}

// EventManagerConfigOpt is a functional option for configuring an EventManager.
type EventManagerConfigOpt func(config *EventManagerConfig)

// Apply applies the given EventManagerConfigOpt(s) to the EventManagerConfig.
func (c *EventManagerConfig) Apply(opts []EventManagerConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithListeners adds the given EventListener(s) to the EventManagerConfig.
func WithListeners(listeners ...EventListener) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.EventListeners = append(config.EventListeners, listeners...)
	}
}

func WithListenerFunc[E Event](listenerFunc func(e E)) EventManagerConfigOpt {
	return WithListeners(ListenerFunc[E](listenerFunc))
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
