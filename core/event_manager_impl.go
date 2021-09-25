package core

import (
	"io"
	"runtime/debug"

	"github.com/DisgoOrg/disgo/json"

	"github.com/DisgoOrg/disgo/discord"
)

var _ EventManager = (*eventManagerImpl)(nil)

func NewEventManager(bot *Bot, listeners []EventListener) EventManager {
	return &eventManagerImpl{
		gatewayEventHandlers:   GetGatewayHandlers(),
		httpServerEventHandler: &httpserverHandlerInteractionCreate{},
		bot:                    bot,
		listeners:              listeners,
	}
}

// eventManagerImpl is the implementation of core.EventManager
type eventManagerImpl struct {
	gatewayEventHandlers   map[discord.GatewayEventType]GatewayEventHandler
	httpServerEventHandler HTTPServerEventHandler
	bot                    *Bot
	listeners              []EventListener
}

// Bot returns the core.Bot instance used by the core.EventManager
func (e *eventManagerImpl) Bot() *Bot {
	return e.bot
}

// Close closes all goroutines created by the core.EventManager
func (e *eventManagerImpl) Close() {
	e.Bot().Logger.Info("closing eventManager goroutines...")
}

// HandleGateway calls the correct core.EventHandler
func (e *eventManagerImpl) HandleGateway(gatewayEventType discord.GatewayEventType, sequenceNumber int, reader io.Reader) {
	if handler, ok := e.gatewayEventHandlers[gatewayEventType]; ok {
		v := handler.New()
		if v != nil {
			if err := json.NewDecoder(reader).Decode(&v); err != nil {
				e.Bot().Logger.Errorf("error while unmarshalling event '%s'. error: %s", gatewayEventType, err.Error())
			}
		}
		handler.HandleGatewayEvent(e.Bot(), sequenceNumber, v)
	} else {
		e.Bot().Logger.Warnf("no handler for gateway event '%s' found", gatewayEventType)
	}
}

// HandleHTTP calls the correct core.EventHandler
func (e *eventManagerImpl) HandleHTTP(responseChannel chan<- discord.InteractionResponse, reader io.Reader) {
	v := e.httpServerEventHandler.New()
	if err := json.NewDecoder(reader).Decode(&v); err != nil {
		e.Bot().Logger.Error("error while unmarshalling httpserver event. error: ", err)
	}
	e.httpServerEventHandler.HandleHTTPEvent(e.Bot(), responseChannel, v)
}

// Dispatch dispatches a new event to the client
func (e *eventManagerImpl) Dispatch(event Event) {
	for i := range e.listeners {
		listener := e.listeners[i]
		go func() {
			defer func() {
				if r := recover(); r != nil {
					e.Bot().Logger.Error("recovered from panic in event listener: ", r)
					debug.PrintStack()
					return
				}
			}()
			listener.OnEvent(event)
		}()
	}
}

// AddEventListeners adds one or more core.EventListener(s) to the core.EventManager
func (e *eventManagerImpl) AddEventListeners(listeners ...EventListener) {
	for _, listener := range listeners {
		e.listeners = append(e.listeners, listener)
	}
}

// RemoveEventListeners removes one or more core.EventListener(s) from the core.EventManager
func (e *eventManagerImpl) RemoveEventListeners(listeners ...EventListener) {
	for _, listener := range listeners {
		for i, l := range e.listeners {
			if l == listener {
				e.listeners = append(e.listeners[:i], e.listeners[i+1:]...)
				break
			}
		}
	}
}
