package internal

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
)

func newEventManagerImpl(disgo api.Disgo, listeners []api.EventListener) api.EventManager {
	eventManager := &EventManagerImpl{
		disgo: disgo,
		channel:   make(chan api.GenericEvent),
		listeners: listeners,
		handlers:  getHandlers(),
	}
	go eventManager.ListenEvents()
	return eventManager
}

type EventManagerImpl struct {
	disgo     api.Disgo
	listeners []api.EventListener
	handlers  *map[string]api.GatewayEventProvider
	channel   chan api.GenericEvent
}

func (e EventManagerImpl) Handle(name string, payload json.RawMessage) {
	if handler, ok := (*e.handlers)[name]; ok {
		eventPayload := handler.New()
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			log.Errorf("error while unmarshaling event. error: %s", err)
		}
		handler.Handle(e.disgo, e, eventPayload)
	}
}

func (e EventManagerImpl) Dispatch(event api.GenericEvent) {
	e.channel <- event
}

func (e EventManagerImpl) AddEventListeners(listeners ...api.EventListener) {
	for _, listener := range listeners {
		e.listeners = append(e.listeners, listener)
	}
}

func (e EventManagerImpl) ListenEvents() {
	/*defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered event listen goroutine error: %s", r)
			debug.PrintStack()
			e.ListenEvents()
			return
		}
		log.Infof("closing event channel...")
		close(e.channel)
	}()*/
	for {
		event := <-e.channel
		for _, listener := range e.listeners {
			listener.OnEvent(event)
		}
	}
}
