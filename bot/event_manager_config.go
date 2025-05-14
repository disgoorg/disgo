package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/gateway"
)

func defaultEventManagerConfig() eventManagerConfig {
	return eventManagerConfig{
		Logger: slog.Default(),
	}
}

type eventManagerConfig struct {
	Logger             *slog.Logger
	EventListeners     []EventListener
	AsyncEventsEnabled bool

	GatewayHandlers   map[gateway.EventType]GatewayEventHandler
	HTTPServerHandler HTTPServerEventHandler
}

// EventManagerConfigOpt is a functional option for configuring an EventManager.
type EventManagerConfigOpt func(config *eventManagerConfig)

func (c *eventManagerConfig) apply(opts []EventManagerConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "bot_event_manager"))
}

// WithEventManagerLogger overrides the default Logger in the eventManagerConfig.
func WithEventManagerLogger(logger *slog.Logger) EventManagerConfigOpt {
	return func(config *eventManagerConfig) {
		config.Logger = logger
	}
}

// WithListeners adds the given EventListener(s) to the eventManagerConfig.
func WithListeners(listeners ...EventListener) EventManagerConfigOpt {
	return func(config *eventManagerConfig) {
		config.EventListeners = append(config.EventListeners, listeners...)
	}
}

// WithListenerFunc adds the given func(e E) to the eventManagerConfig.
func WithListenerFunc[E Event](f func(e E)) EventManagerConfigOpt {
	return WithListeners(NewListenerFunc(f))
}

// WithListenerChan adds the given chan<- E to the eventManagerConfig.
func WithListenerChan[E Event](c chan<- E) EventManagerConfigOpt {
	return WithListeners(NewListenerChan(c))
}

// WithAsyncEventsEnabled enables/disables the async events.
func WithAsyncEventsEnabled() EventManagerConfigOpt {
	return func(config *eventManagerConfig) {
		config.AsyncEventsEnabled = true
	}
}

// WithGatewayHandlers overrides the default GatewayEventHandler(s) in the eventManagerConfig.
func WithGatewayHandlers(handlers map[gateway.EventType]GatewayEventHandler) EventManagerConfigOpt {
	return func(config *eventManagerConfig) {
		config.GatewayHandlers = handlers
	}
}

// WithHTTPServerHandler overrides the given HTTPServerEventHandler in the eventManagerConfig.
func WithHTTPServerHandler(handler HTTPServerEventHandler) EventManagerConfigOpt {
	return func(config *eventManagerConfig) {
		config.HTTPServerHandler = handler
	}
}
