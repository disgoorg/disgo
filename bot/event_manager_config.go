package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/gateway"
)

// DefaultEventManagerConfig returns a new EventManagerConfig with all default values.
func DefaultEventManagerConfig() *EventManagerConfig {
	return &EventManagerConfig{
		Logger: slog.Default(),
	}
}

// EventManagerConfig can be used to configure the EventManager.
type EventManagerConfig struct {
	Logger             *slog.Logger
	EventListeners     []EventListener
	AsyncEventsEnabled bool

	GatewayHandlers   map[gateway.EventType]GatewayEventHandler
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

// WithEventManagerLogger overrides the default logger in the EventManagerConfig.
func WithEventManagerLogger(logger *slog.Logger) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.Logger = logger
	}
}

// WithListeners adds the given EventListener(s) to the EventManagerConfig.
func WithListeners(listeners ...EventListener) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.EventListeners = append(config.EventListeners, listeners...)
	}
}

// WithListenerFunc adds the given func(e E) to the EventManagerConfig.
func WithListenerFunc[E Event](f func(e E)) EventManagerConfigOpt {
	return WithListeners(NewListenerFunc(f))
}

// WithListenerChan adds the given chan<- E to the EventManagerConfig.
func WithListenerChan[E Event](c chan<- E) EventManagerConfigOpt {
	return WithListeners(NewListenerChan(c))
}

// WithAsyncEventsEnabled enables/disables the async events.
func WithAsyncEventsEnabled() EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.AsyncEventsEnabled = true
	}
}

// WithGatewayHandlers overrides the default GatewayEventHandler(s) in the EventManagerConfig.
func WithGatewayHandlers(handlers map[gateway.EventType]GatewayEventHandler) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.GatewayHandlers = handlers
	}
}

// WithHTTPServerHandler overrides the given HTTPServerEventHandler in the EventManagerConfig.
func WithHTTPServerHandler(handler HTTPServerEventHandler) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.HTTPServerHandler = handler
	}
}
