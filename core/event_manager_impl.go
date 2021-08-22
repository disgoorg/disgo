package core

import (
	"encoding/json"
	"io"
	"runtime/debug"

	"github.com/DisgoOrg/disgo/discord"
)

var (
	GatewayEventHandlers = map[discord.GatewayEventType]GatewayEventHandler{}
	HTTPEventHandlers    = map[discord.GatewayEventType]HTTPEventHandler{}
)

var _ EventManager = (*EventManagerImpl)(nil)

func NewEventManager(disgo Disgo, listeners []EventListener) EventManager {
	eventManager := &EventManagerImpl{
		disgo:     disgo,
		channel:   make(chan Event),
		listeners: listeners,
	}

	go eventManager.ListenEvents()
	return eventManager
}

// EventManagerImpl is the implementation of api.EventManager
type EventManagerImpl struct {
	disgo     Disgo
	listeners []EventListener
	channel   chan Event
}

// Disgo returns the api.Disgo instance used by the api.EventManager
func (e *EventManagerImpl) Disgo() Disgo {
	return e.disgo
}

// Close closes all goroutines created by the api.EventManager
func (e *EventManagerImpl) Close() {
	e.Disgo().Logger().Info("closing eventManager goroutines...")
	close(e.channel)
}

// HandleGateway calls the correct api.EventHandler
func (e *EventManagerImpl) HandleGateway(gatewayEventType discord.GatewayEventType, sequenceNumber int, payload io.Reader) {
	if handler, ok := GatewayEventHandlers[gatewayEventType]; ok {
		eventPayload := handler.New()
		if err := json.NewDecoder(payload).Decode(&eventPayload); err != nil {
			e.disgo.Logger().Errorf("error while unmarshalling event. error: %s", err)
		}
		handler.HandleGatewayEvent(e.disgo, e, sequenceNumber, eventPayload)
	}
}

// HandleHTTP calls the correct api.EventHandler
func (e *EventManagerImpl) HandleHTTP(gatewayEventType discord.GatewayEventType, c chan discord.InteractionResponse, payload io.Reader) {
	if handler, ok := HTTPEventHandlers[gatewayEventType]; ok {
		eventPayload := handler.New()
		if err := json.NewDecoder(payload).Decode(&eventPayload); err != nil {
			e.disgo.Logger().Errorf("error while unmarshalling httpserver event. error: %s", err)
		}
		handler.HandleHTTPEvent(e.disgo, e, c, eventPayload)
	}
}

// Dispatch dispatches a new event to the client
func (e *EventManagerImpl) Dispatch(event Event) {
	go func() {
		e.channel <- event
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

// ListenEvents starts the event goroutine
func (e *EventManagerImpl) ListenEvents() {
	defer func() {
		if r := recover(); r != nil {
			e.Disgo().Logger().Panicf("recovered event listen goroutine error: %s", r)
			debug.PrintStack()
			e.ListenEvents()
			return
		}
		e.Disgo().Logger().Infof("closed event goroutine")
	}()
	for {
		event, ok := <-e.channel
		if !ok {
			return
		}
		for _, listener := range e.listeners {
			listener.OnEvent(event)
		}
	}
}
