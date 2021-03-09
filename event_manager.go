package disgo

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type GatewayEventProvider interface {
	New() interface{}
	Handle(EventManager, interface{})
}

type EventListener interface {
	OnEvent(interface{})
}

type EventManager interface {
	AddEventListeners(...EventListener)
	handle(string, json.RawMessage)
	Dispatch(GenericEvent)
	Disgo() Disgo
}

type EventManagerImpl struct {
	disgo     Disgo
	listeners *[]*EventListener
	handlers  *map[string]GatewayEventProvider
	channel   chan GenericEvent
}

func (e EventManagerImpl) Disgo() Disgo {
	return e.disgo
}

func (e EventManagerImpl) handle(name string, payload json.RawMessage) {
	if handler, ok := (*e.handlers)[name]; ok {
		eventPayload := handler.New()
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			log.Errorf("error while unmarshaling event. error: %s", err)
		}
		handler.Handle(e, eventPayload)
	}
}

func (e EventManagerImpl) Dispatch(event GenericEvent) {
	e.channel <- event
}

func (e EventManagerImpl) AddEventListeners(listeners ...EventListener) {
	for _, listener := range listeners {
		*e.listeners = append(*e.listeners, &listener)
	}
}

func (e EventManagerImpl) listenEvents() {
	defer func() {
		log.Infof("closing event channel...")
		close(e.channel)
	}()
	for {
		event := <-e.channel
		for _, listener := range *e.listeners {
			(*listener).OnEvent(event)
		}
	}
}

func newEventManager(disgo Disgo) EventManager {
	manager := EventManagerImpl{
		disgo:     disgo,
		channel:   make(chan GenericEvent),
		listeners: &[]*EventListener{},
		handlers:  GetHandlers(),
	}
	go manager.listenEvents()
	return manager
}

func (e EventManagerImpl) AddHandler(event string, handler GatewayEventProvider) {
	(*e.handlers)[event] = handler
}
