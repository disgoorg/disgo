package bot

import (
	"runtime/debug"
	"sync"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
)

var _ EventManager = (*eventManagerImpl)(nil)

// NewEventManager returns a new EventManager with the EventManagerConfigOpt(s) applied.
func NewEventManager(client Client, opts ...EventManagerConfigOpt) EventManager {
	config := DefaultEventManagerConfig()
	config.Apply(opts)

	return &eventManagerImpl{
		client: client,
		config: *config,
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
	HandleHTTPEvent(respondFunc httpserver.RespondFunc, event gateway.EventInteractionCreate)

	// DispatchEvent dispatches a new Event to the Client's EventListener(s)
	DispatchEvent(event Event)
}

// EventListener is used to create new EventListener to listen to events
type EventListener interface {
	OnEvent(event Event)
}

var _ EventListener = (*ListenerFunc[Event])(nil)

// NewListenerFunc returns a new ListenerFunc for the given func(e E)
func NewListenerFunc[E Event](f func(e E)) *ListenerFunc[E] {
	return &ListenerFunc[E]{F: f}
}

// ListenerFunc is a wrapper for a func(e E) as functions are not comparable
type ListenerFunc[E Event] struct {
	F func(e E)
}

// OnEvent calls the func(e E) if E is Event
func (l *ListenerFunc[E]) OnEvent(e Event) {
	if event, ok := e.(E); ok {
		l.F(event)
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
	HandleHTTPEvent(client Client, respondFunc httpserver.RespondFunc, event gateway.EventInteractionCreate)
}

type eventManagerImpl struct {
	client          Client
	eventListenerMu sync.Mutex
	config          EventManagerConfig

	mu sync.Mutex
}

func (e *eventManagerImpl) HandleGatewayEvent(gatewayEventType gateway.EventType, sequenceNumber int, shardID int, event gateway.EventData) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if handler, ok := e.config.GatewayHandlers[gatewayEventType]; ok {
		handler.HandleGatewayEvent(e.client, sequenceNumber, shardID, event)
	} else {
		e.client.Logger().Warnf("no handler for gateway event '%s' found", gatewayEventType)
	}
}

func (e *eventManagerImpl) HandleHTTPEvent(respondFunc httpserver.RespondFunc, event gateway.EventInteractionCreate) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.config.HTTPServerHandler.HandleHTTPEvent(e.client, respondFunc, event)
}

func (e *eventManagerImpl) DispatchEvent(event Event) {
	defer func() {
		if r := recover(); r != nil {
			e.client.Logger().Errorf("recovered from panic in event listener: %+v\nstack: %s", r, string(debug.Stack()))
			return
		}
	}()
	e.eventListenerMu.Lock()
	defer e.eventListenerMu.Unlock()
	for i := range e.config.EventListeners {
		if e.config.AsyncEventsEnabled {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						e.client.Logger().Errorf("recovered from panic in event listener: %+v\nstack: %s", r, string(debug.Stack()))
						return
					}
				}()
				e.config.EventListeners[i].OnEvent(event)
			}()
			continue
		}
		e.config.EventListeners[i].OnEvent(event)
	}
}

func (e *eventManagerImpl) AddEventListeners(listeners ...EventListener) {
	e.eventListenerMu.Lock()
	defer e.eventListenerMu.Unlock()
	e.config.EventListeners = append(e.config.EventListeners, listeners...)
}

func (e *eventManagerImpl) RemoveEventListeners(listeners ...EventListener) {
	e.eventListenerMu.Lock()
	defer e.eventListenerMu.Unlock()
	for _, listener := range listeners {
		for i, l := range e.config.EventListeners {
			if l == listener {
				e.config.EventListeners = append(e.config.EventListeners[:i], e.config.EventListeners[i+1:]...)
				break
			}
		}
	}
}
