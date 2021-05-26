package internal

import (
	"encoding/json"
	"runtime/debug"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/internal/handlers"
)

func newEventManagerImpl(disgo api.Disgo, listeners []api.EventListener) api.EventManager {
	eventManager := &EventManagerImpl{
		disgo:     disgo,
		channel:   make(chan api.Event),
		listeners: listeners,
		handlers:  map[api.GatewayEventType]api.EventHandler{},
	}
	for _, handler := range handlers.GetAllHandlers() {
		eventManager.handlers[handler.Event()] = handler
	}
	go eventManager.ListenEvents()
	return eventManager
}

// EventManagerImpl is the implementation of api.EventManager
type EventManagerImpl struct {
	disgo     api.Disgo
	listeners []api.EventListener
	handlers  map[api.GatewayEventType]api.EventHandler
	channel   chan api.Event
}

// Disgo returns the api.Disgo instance used by the api.EventManager
func (e *EventManagerImpl) Disgo() api.Disgo {
	return e.disgo
}

// Close closes all goroutines created by the api.EventManager
func (e *EventManagerImpl) Close() {
	e.Disgo().Logger().Info("closing eventManager goroutines...")
	close(e.channel)
}

// Handle calls the correct api.EventHandler
func (e *EventManagerImpl) Handle(name api.GatewayEventType, c chan *api.InteractionResponse, sequenceNumber int, payload json.RawMessage) {
	if handler, ok := e.handlers[name]; ok {
		eventPayload := handler.New()
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			e.disgo.Logger().Errorf("error while unmarshalling event. error: %s", err)
		}

		switch h := handler.(type) {
		case api.GatewayEventHandler:
			h.HandleGatewayEvent(e.disgo, e, sequenceNumber, eventPayload)
		case api.WebhookEventHandler:
			h.HandleWebhookEvent(e.disgo, e, c, eventPayload)
		default:
			e.Disgo().Logger().Errorf("no event handler found for: %s", name)
		}
	}
}

// Dispatch dispatches a new event to the client
func (e *EventManagerImpl) Dispatch(event api.Event) {
	go func() {
		e.channel <- event
	}()
}

// AddEventListeners adds one or more api.EventListener(s) to the api.EventManager
func (e *EventManagerImpl) AddEventListeners(listeners ...api.EventListener) {
	for _, listener := range listeners {
		e.listeners = append(e.listeners, listener)
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
