package bot

import (
	"log/slog"
	"runtime/debug"
	"sync"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpgateway"
)

// EventListener is used to create new EventListener to listen to events
type EventListener interface {
	OnEvent(event Event)
}

// NewListenerFunc returns a new EventListener for the given func(e E)
func NewListenerFunc[E Event](f func(e E)) EventListener {
	return &listenerFunc[E]{f: f}
}

type listenerFunc[E Event] struct {
	f func(e E)
}

func (l *listenerFunc[E]) OnEvent(e Event) {
	if event, ok := e.(E); ok {
		l.f(event)
	}
}

// NewListenerChan returns a new EventListener for the given chan<- Event
func NewListenerChan[E Event](c chan<- E) EventListener {
	return &listenerChan[E]{c: c}
}

type listenerChan[E Event] struct {
	c chan<- E
}

func (l *listenerChan[E]) OnEvent(e Event) {
	if event, ok := e.(E); ok {
		l.c <- event
	}
}

// Event the basic interface each event implement
type Event interface {
	Client() *Client
	SequenceNumber() int
}

// GatewayEventHandler is used to handle Gateway Event(s)
type GatewayEventHandler interface {
	HandleGatewayEvent(client *Client, message gateway.Message, shardID int)
}

// HTTPInteractionEventHandler is used to handle HTTP Event(s)
type HTTPInteractionEventHandler interface {
	HandleHTTPInteraction(client *Client, respond httpgateway.RespondFunc, event httpgateway.EventInteractionCreate)
}

type HTTPGatewayEventHandler interface {
	HandleHTTPGatewayEvent(client *Client, ack func(), message httpgateway.Message)
}

// EventManager lets you listen for specific events triggered by raw Gateway events
type EventManager interface {
	// AddEventListeners adds one or more EventListener(s) to the EventManager
	AddEventListeners(eventListeners ...EventListener)

	// RemoveEventListeners removes one or more EventListener(s) from the EventManager
	RemoveEventListeners(eventListeners ...EventListener)

	// HandleGatewayEvent calls the correct GatewayEventHandler for the payload
	HandleGatewayEvent(gateway gateway.Gateway, message gateway.Message)

	// HandleHTTPInteractionEvent calls the HTTPInteractionEventHandler for the payload
	HandleHTTPInteractionEvent(respond httpgateway.RespondFunc, event httpgateway.EventInteractionCreate)

	// HandleHTTPGatewayEvent calls the HTTPInteractionEventHandler for the payload
	HandleHTTPGatewayEvent(ack func(), message httpgateway.Message)

	// DispatchEvent dispatches a new Event to the Client's EventListener(s)
	DispatchEvent(event Event)
}

var _ EventManager = (*eventManagerImpl)(nil)

// NewEventManager returns a new EventManager with the EventManagerConfigOpt(s) applied.
func NewEventManager(client *Client, opts ...EventManagerConfigOpt) EventManager {
	cfg := defaultEventManagerConfig()
	cfg.apply(opts)

	return &eventManagerImpl{
		client:                 client,
		logger:                 cfg.Logger,
		asyncEventsEnabled:     cfg.AsyncEventsEnabled,
		eventListeners:         cfg.EventListeners,
		gatewayHandler:         cfg.GatewayHandler,
		httpInteractionHandler: cfg.HTTPInteractionHandler,
		httpGatewayHandler:     cfg.HTTPGatewayHandler,
	}
}

type eventManagerImpl struct {
	client *Client
	logger *slog.Logger

	asyncEventsEnabled bool
	eventListenerMu    sync.Mutex
	eventListeners     []EventListener

	mu                     sync.Mutex
	gatewayHandler         GatewayEventHandler
	httpInteractionHandler HTTPInteractionEventHandler
	httpGatewayHandler     HTTPGatewayEventHandler
}

func (e *eventManagerImpl) HandleGatewayEvent(gateway gateway.Gateway, message gateway.Message) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.gatewayHandler.HandleGatewayEvent(e.client, message, gateway.ShardID())
}

func (e *eventManagerImpl) HandleHTTPInteractionEvent(respond httpgateway.RespondFunc, event httpgateway.EventInteractionCreate) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.httpInteractionHandler.HandleHTTPInteraction(e.client, respond, event)
}

func (e *eventManagerImpl) HandleHTTPGatewayEvent(ack func(), message httpgateway.Message) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.httpGatewayHandler.HandleHTTPGatewayEvent(e.client, ack, message)
}

func (e *eventManagerImpl) DispatchEvent(event Event) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Error("recovered from panic in event listener", slog.Any("arg", r), slog.String("stack", string(debug.Stack())))
			return
		}
	}()
	e.eventListenerMu.Lock()
	defer e.eventListenerMu.Unlock()
	for _, listener := range e.eventListeners {
		if e.asyncEventsEnabled {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						e.logger.Error("recovered from panic in event listener", slog.Any("arg", r), slog.String("stack", string(debug.Stack())))
						return
					}
				}()
				listener.OnEvent(event)
			}()
			continue
		}
		listener.OnEvent(event)
	}
}

func (e *eventManagerImpl) AddEventListeners(listeners ...EventListener) {
	e.eventListenerMu.Lock()
	defer e.eventListenerMu.Unlock()
	e.eventListeners = append(e.eventListeners, listeners...)
}

func (e *eventManagerImpl) RemoveEventListeners(listeners ...EventListener) {
	e.eventListenerMu.Lock()
	defer e.eventListenerMu.Unlock()
	for _, listener := range listeners {
		for i, l := range e.eventListeners {
			if l == listener {
				e.eventListeners = append(e.eventListeners[:i], e.eventListeners[i+1:]...)
				break
			}
		}
	}
}
