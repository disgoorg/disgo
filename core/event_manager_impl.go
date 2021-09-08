package core

import (
	"io"
	"runtime/debug"

	"github.com/DisgoOrg/disgo/json"

	"github.com/DisgoOrg/disgo/discord"
)

var (
	GatewayEventHandlers   = map[discord.GatewayEventType]GatewayEventHandler{}
	HTTPServerEventHandler HTTPEventHandler
)

var _ EventManager = (*EventManagerImpl)(nil)

func NewEventManager(bot *Bot, listeners []EventListener) EventManager {
	return &EventManagerImpl{
		bot:       bot,
		listeners: listeners,
	}
}

// EventManagerImpl is the implementation of api.EventManager
type EventManagerImpl struct {
	bot       *Bot
	listeners []EventListener
}

// Bot returns the api.Bot instance used by the api.EventManager
func (e *EventManagerImpl) Bot() *Bot {
	return e.bot
}

// Close closes all goroutines created by the api.EventManager
func (e *EventManagerImpl) Close() {
	e.Bot().Logger.Info("closing eventManager goroutines...")
}

// HandleGateway calls the correct api.EventHandler
func (e *EventManagerImpl) HandleGateway(gatewayEventType discord.GatewayEventType, sequenceNumber int, reader io.Reader) {
	println("handling event")
	if handler, ok := GatewayEventHandlers[gatewayEventType]; ok {
		v := handler.New()
		if err := json.NewDecoder(reader).Decode(&v); err != nil {
			e.Bot().Logger.Error("error while unmarshalling event. error: ", err)
		}
		handler.HandleGatewayEvent(e.Bot(), sequenceNumber, v)
	} else {
		e.Bot().Logger.Warnf("no handler for gateway event '%s' found", gatewayEventType)
	}
}

// HandleHTTP calls the correct api.EventHandler
func (e *EventManagerImpl) HandleHTTP(c chan discord.InteractionResponse, reader io.Reader) {
	v := HTTPServerEventHandler.New()
	if err := json.NewDecoder(reader).Decode(&v); err != nil {
		e.Bot().Logger.Error("error while unmarshalling httpserver event. error: ", err)
	}
	HTTPServerEventHandler.HandleHTTPEvent(e.Bot(), c, v)
}

// Dispatch dispatches a new event to the client
func (e *EventManagerImpl) Dispatch(event Event) {
	println("called")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				e.Bot().Logger.Panicf("recovered from listener panic error: %s", r)
				debug.PrintStack()
				return
			}
		}()
		for i, listener := range e.listeners {
			println("listener index: ", i)
			listener.OnEvent(event)
		}
	}()
}

// AddEventListeners adds one or more api.EventListener(s) to the api.EventManager
func (e *EventManagerImpl) AddEventListeners(listeners ...EventListener) {
	for _, listener := range listeners {
		e.listeners = append(e.listeners, listener)
	}
}

// RemoveEventListeners removes one or more api.EventListener(s) from the api.EventManager
func (e *EventManagerImpl) RemoveEventListeners(listeners ...EventListener) {
	for _, listener := range listeners {
		for i, l := range e.listeners {
			if l == listener {
				e.listeners = append(e.listeners[:i], e.listeners[i+1:]...)
				break
			}
		}
	}
}
