package internal

import (
	"encoding/json"
	"runtime/debug"

	log "github.com/sirupsen/logrus"

	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/internal/handlers"
)

func newEventManagerImpl(disgo api.Disgo, listeners []api.EventListener) api.EventManager {
	eventManager := &EventManagerImpl{
		disgo:     disgo,
		channel:   make(chan api.Event),
		listeners: listeners,
		handlers:  map[api.GatewayEventName]api.EventHandler{},
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
	handlers  map[api.GatewayEventName]api.EventHandler
	channel   chan api.Event
}

func (e EventManagerImpl) Close() {
	log.Info("closing eventManager goroutines...")
	close(e.channel)
}

// Handle calls the correct api.EventHandler
func (e EventManagerImpl) Handle(name api.GatewayEventName, payload json.RawMessage, c chan interface{}) {
	if handler, ok := e.handlers[name]; ok {
		eventPayload := handler.New()
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			log.Errorf("error while unmarshaling event. error: %s", err)
		}
		switch h := handler.(type) {
		case api.GatewayEventHandler:
			h.Handle(e.disgo, e, eventPayload)
		case api.WebhookEventHandler:
			h.Handle(e.disgo, e, c, eventPayload)
		}
	}
}

// Dispatch dispatches a new event to the client
func (e EventManagerImpl) Dispatch(event api.Event) {
	e.channel <- event
}

// AddEventListeners adds one or more api.EventListener(s) to the api.EventManager
func (e EventManagerImpl) AddEventListeners(listeners ...api.EventListener) {
	for _, listener := range listeners {
		e.listeners = append(e.listeners, listener)
	}
}

// ListenEvents starts the event goroutine
func (e EventManagerImpl) ListenEvents() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered event listen goroutine error: %s", r)
			debug.PrintStack()
			e.ListenEvents()
			return
		}
		log.Infof("closed event goroutine")
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
