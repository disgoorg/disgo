package bot

import (
	"io"
	"runtime/debug"
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/json"
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
	// RawEventsEnabled returns whether events.RawEvent are enabled
	RawEventsEnabled() bool

	// AddEventListeners adds one or more EventListener(s) to the EventManager
	AddEventListeners(eventListeners ...EventListener)

	// RemoveEventListeners removes one or more EventListener(s) from the EventManager
	RemoveEventListeners(eventListeners ...EventListener)

	// HandleGatewayEvent calls the correct GatewayEventHandler for the payload
	HandleGatewayEvent(gatewayEventType discord.GatewayEventType, sequenceNumber int, shardID int, payload io.Reader)

	// HandleHTTPEvent calls the HTTPServerEventHandler for the payload
	HandleHTTPEvent(respondFunc httpserver.RespondFunc, payload io.Reader)

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
	EventType() discord.GatewayEventType
	New() any
	HandleGatewayEvent(client Client, sequenceNumber int, shardID int, v any)
}

// HTTPServerEventHandler is used to handle HTTP Event(s)
type HTTPServerEventHandler interface {
	New() any
	HandleHTTPEvent(client Client, respondFunc func(response discord.InteractionResponse) error, v any)
}

type eventManagerImpl struct {
	client          Client
	eventListenerMu sync.Mutex
	config          EventManagerConfig

	mu sync.Mutex
}

func (e *eventManagerImpl) RawEventsEnabled() bool {
	return e.config.RawEventsEnabled
}

func (e *eventManagerImpl) HandleGatewayEvent(gatewayEventType discord.GatewayEventType, sequenceNumber int, shardID int, reader io.Reader) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if handler, ok := e.config.GatewayHandlers[gatewayEventType]; ok {
		v := handler.New()
		if v != nil {
			if err := json.NewDecoder(reader).Decode(&v); err != nil {
				e.client.Logger().Errorf("error while unmarshalling event '%s'. error: %s", gatewayEventType, err.Error())
				return
			}
		}
		handler.HandleGatewayEvent(e.client, sequenceNumber, shardID, v)
	} else {
		e.client.Logger().Warnf("no handler for gateway event '%s' found", gatewayEventType)
	}
}

func (e *eventManagerImpl) HandleHTTPEvent(respondFunc httpserver.RespondFunc, reader io.Reader) {
	e.mu.Lock()
	defer e.mu.Unlock()
	v := e.config.HTTPServerHandler.New()
	if err := json.NewDecoder(reader).Decode(&v); err != nil {
		e.client.Logger().Error("error while unmarshalling httpserver event. error: ", err)
	}
	e.config.HTTPServerHandler.HandleHTTPEvent(e.client, respondFunc, v)
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
