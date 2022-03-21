package bot

import (
	"io"
	"runtime/debug"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/json"
)

var _ EventManager = (*eventManagerImpl)(nil)

func NewEventManager(client Client, opts ...EventManagerConfigOpt) EventManager {
	config := &DefaultEventManagerConfig
	config.Apply(opts)

	return &eventManagerImpl{
		bot:    client,
		config: *config,
	}
}

// EventManager lets you listen for specific events triggered by raw BotGateway events
type EventManager interface {
	Client() Client
	Config() EventManagerConfig

	AddEventListeners(eventListeners ...EventListener)
	RemoveEventListeners(eventListeners ...EventListener)
	HandleGateway(gatewayEventType discord.GatewayEventType, sequenceNumber discord.GatewaySequence, payload io.Reader)
	HandleHTTP(responseChannel chan<- discord.InteractionResponse, payload io.Reader)
	Dispatch(event Event)
}

// EventListener is used to create new EventListener to listen to events
type EventListener interface {
	OnEvent(event Event)
}

// Event the basic interface each event implement
type Event interface {
	Client() Client
	SequenceNumber() discord.GatewaySequence
}

// GatewayEventHandler is used to handle Gateway Event(s)
type GatewayEventHandler interface {
	EventType() discord.GatewayEventType
	New() any
	HandleGatewayEvent(client Client, sequenceNumber discord.GatewaySequence, v any)
}

// HTTPServerEventHandler is used to handle HTTP Event(s)
type HTTPServerEventHandler interface {
	New() any
	HandleHTTPEvent(client Client, responseChannel chan<- discord.InteractionResponse, v any)
}

// eventManagerImpl is the implementation of core.EventManager
type eventManagerImpl struct {
	bot    Client
	config EventManagerConfig
}

// Client returns the core.Client instance used by the core.EventManager
func (e *eventManagerImpl) Client() Client {
	return e.bot
}

func (e *eventManagerImpl) Config() EventManagerConfig {
	return e.config
}

// HandleGateway calls the correct core.EventHandler
func (e *eventManagerImpl) HandleGateway(gatewayEventType discord.GatewayEventType, sequenceNumber discord.GatewaySequence, reader io.Reader) {
	if handler, ok := e.config.GatewayHandlers[gatewayEventType]; ok {
		v := handler.New()
		if v != nil {
			if err := json.NewDecoder(reader).Decode(&v); err != nil {
				e.Client().Logger().Errorf("error while unmarshalling event '%s'. error: %s", gatewayEventType, err.Error())
				return
			}
		}
		handler.HandleGatewayEvent(e.Client(), sequenceNumber, v)
	} else {
		e.Client().Logger().Warnf("no handler for BotGateway event '%s' found", gatewayEventType)
	}
}

// HandleHTTP calls the correct core.EventHandler
func (e *eventManagerImpl) HandleHTTP(responseChannel chan<- discord.InteractionResponse, reader io.Reader) {
	v := e.config.HTTPServerHandler.New()
	if err := json.NewDecoder(reader).Decode(&v); err != nil {
		e.Client().Logger().Error("error while unmarshalling httpserver event. error: ", err)
	}
	e.config.HTTPServerHandler.HandleHTTPEvent(e.Client(), responseChannel, v)
}

// Dispatch dispatches a new event to the client
func (e *eventManagerImpl) Dispatch(event Event) {
	defer func() {
		if r := recover(); r != nil {
			e.Client().Logger().Errorf("recovered from panic in event listener: %+v\nstack: %s", r, string(debug.Stack()))
			return
		}
	}()
	for i := range e.config.EventListeners {
		if e.Config().AsyncEventsEnabled {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						e.Client().Logger().Errorf("recovered from panic in event listener: %+v\nstack: %s", r, string(debug.Stack()))
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
	e.config.EventListeners = append(e.config.EventListeners, listeners...)
}

// RemoveEventListeners removes one or more core.EventListener(s) from the core.EventManager
func (e *eventManagerImpl) RemoveEventListeners(listeners ...EventListener) {
	for _, listener := range listeners {
		for i, l := range e.config.EventListeners {
			if l == listener {
				e.config.EventListeners = append(e.config.EventListeners[:i], e.config.EventListeners[i+1:]...)
				break
			}
		}
	}
}
