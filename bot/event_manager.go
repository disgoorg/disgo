package bot

import (
	"log/slog"
	"runtime/debug"
	"sync"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
)

var _ EventManager = (*eventManagerImpl)(nil)

// NewEventManager returns a new EventManager with the EventManagerConfigOpt(s) applied.
func NewEventManager(client Client, opts ...EventManagerConfigOpt) EventManager {
	cfg := DefaultEventManagerConfig()
	cfg.Apply(opts)
	cfg.Logger = cfg.Logger.With(slog.String("name", "bot_event_manager"))

	return &eventManagerImpl{
		client:             client,
		logger:             cfg.Logger,
		eventListeners:     cfg.EventListeners,
		asyncEventsEnabled: cfg.AsyncEventsEnabled,
		gatewayHandlers:    cfg.GatewayHandlers,
		httpServerHandler:  cfg.HTTPServerHandler,
	}
}

// EventManager lets you listen for specific events triggered by raw gateway events
type EventManager interface {
	// AddEventListeners adds one or more EventListener(s) to the EventManager
	AddEventListeners(eventListeners ...EventListener)

	// RemoveEventListeners removes one or more EventListener(s) from the EventManager
	RemoveEventListeners(eventListeners ...EventListener)

	// HandleGatewayEvent calls the correct GatewayEventHandler for the payload
	HandleGatewayEvent(gatewayEventType gateway.EventType, sequenceNumber int, shardID int, event gateway.EventData)

	// HandleHTTPEvent calls the HTTPServerEventHandler for the payload
	HandleHTTPEvent(respondFunc httpserver.RespondFunc, event httpserver.EventInteractionCreate)

	// DispatchEvent dispatches a new Event to the Client's EventListener(s)
	DispatchEvent(event Event)
}

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
	Client() Client
	SequenceNumber() int
}

// GatewayEventHandler is used to handle Gateway Event(s)
type GatewayEventHandler interface {
	EventType() gateway.EventType
	HandleGatewayEvent(client Client, sequenceNumber int, shardID int, event gateway.EventData)
}

// NewGatewayEventHandler returns a new GatewayEventHandler for the given GatewayEventType and handler func
func NewGatewayEventHandler[T gateway.EventData](eventType gateway.EventType, handleFunc func(client Client, sequenceNumber int, shardID int, event T)) GatewayEventHandler {
	return &genericGatewayEventHandler[T]{eventType: eventType, handleFunc: handleFunc}
}

type genericGatewayEventHandler[T gateway.EventData] struct {
	eventType  gateway.EventType
	handleFunc func(client Client, sequenceNumber int, shardID int, event T)
}

func (h *genericGatewayEventHandler[T]) EventType() gateway.EventType {
	return h.eventType
}

func (h *genericGatewayEventHandler[T]) HandleGatewayEvent(client Client, sequenceNumber int, shardID int, event gateway.EventData) {
	if e, ok := event.(T); ok {
		h.handleFunc(client, sequenceNumber, shardID, e)
	}
}

// HTTPServerEventHandler is used to handle HTTP Event(s)
type HTTPServerEventHandler interface {
	HandleHTTPEvent(client Client, respondFunc httpserver.RespondFunc, event httpserver.EventInteractionCreate)
}

type eventManagerImpl struct {
	mu sync.Mutex

	client             Client
	logger             *slog.Logger
	eventListenerMu    sync.Mutex
	eventListeners     []EventListener
	asyncEventsEnabled bool
	gatewayHandlers    map[gateway.EventType]GatewayEventHandler
	httpServerHandler  HTTPServerEventHandler
}

func (e *eventManagerImpl) HandleGatewayEvent(gatewayEventType gateway.EventType, sequenceNumber int, shardID int, event gateway.EventData) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if handler, ok := e.gatewayHandlers[gatewayEventType]; ok {
		handler.HandleGatewayEvent(e.client, sequenceNumber, shardID, event)
	} else {
		e.logger.Warn("no handler for gateway event found", slog.Any("event_type", gatewayEventType))
	}
}

func (e *eventManagerImpl) HandleHTTPEvent(respondFunc httpserver.RespondFunc, event httpserver.EventInteractionCreate) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.httpServerHandler.HandleHTTPEvent(e.client, respondFunc, event)
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
	for i := range e.eventListeners {
		if e.asyncEventsEnabled {
			go func(i int) {
				defer func() {
					if r := recover(); r != nil {
						e.logger.Error("recovered from panic in event listener", slog.Any("arg", r), slog.String("stack", string(debug.Stack())))
						return
					}
				}()
				e.eventListeners[i].OnEvent(event)
			}(i)
			continue
		}
		e.eventListeners[i].OnEvent(event)
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
