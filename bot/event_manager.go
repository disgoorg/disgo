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
	RawEventsEnabled() bool

	AddEventListeners(eventListeners ...EventListener)
	RemoveEventListeners(eventListeners ...EventListener)

	HandleGatewayEvent(gatewayEventType discord.GatewayEventType, sequenceNumber int, payload io.Reader)
	HandleHTTPEvent(respondFunc httpserver.RespondFunc, payload io.Reader)
	DispatchEvent(event Event)
}

// EventListener is used to create new EventListener to listen to events
type EventListener interface {
	OnEvent(event Event)
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
	HandleGatewayEvent(client Client, sequenceNumber int, v any)
}

// HTTPServerEventHandler is used to handle HTTP Event(s)
type HTTPServerEventHandler interface {
	New() any
	HandleHTTPEvent(client Client, respondFunc func(response discord.InteractionResponse) error, v any)
}

// eventManagerImpl is the implementation of core.EventManager
type eventManagerImpl struct {
	client          Client
	eventListenerMu sync.Mutex
	config          EventManagerConfig
}

func (e *eventManagerImpl) RawEventsEnabled() bool {
	return e.config.RawEventsEnabled
}

// HandleGatewayEvent calls the correct core.EventHandler
func (e *eventManagerImpl) HandleGatewayEvent(gatewayEventType discord.GatewayEventType, sequenceNumber int, reader io.Reader) {
	if handler, ok := e.config.GatewayHandlers[gatewayEventType]; ok {
		v := handler.New()
		if v != nil {
			if err := json.NewDecoder(reader).Decode(&v); err != nil {
				e.client.Logger().Errorf("error while unmarshalling event '%s'. error: %s", gatewayEventType, err.Error())
				return
			}
		}
		handler.HandleGatewayEvent(e.client, sequenceNumber, v)
	} else {
		e.client.Logger().Warnf("no handler for gateway event '%s' found", gatewayEventType)
	}
}

// HandleHTTPEvent calls the correct core.EventHandler
func (e *eventManagerImpl) HandleHTTPEvent(respondFunc httpserver.RespondFunc, reader io.Reader) {
	v := e.config.HTTPServerHandler.New()
	if err := json.NewDecoder(reader).Decode(&v); err != nil {
		e.client.Logger().Error("error while unmarshalling httpserver event. error: ", err)
	}
	e.config.HTTPServerHandler.HandleHTTPEvent(e.client, respondFunc, v)
}

// DispatchEvent dispatches a new event to the client
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

// AddEventListeners adds one or more core.EventListener(s) to the core.EventManager
func (e *eventManagerImpl) AddEventListeners(listeners ...EventListener) {
	e.eventListenerMu.Lock()
	defer e.eventListenerMu.Unlock()
	e.config.EventListeners = append(e.config.EventListeners, listeners...)
}

// RemoveEventListeners removes one or more core.EventListener(s) from the core.EventManager
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
