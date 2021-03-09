package internal

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/DiscoOrg/disgo/api"
)

type EventManagerImpl struct {
	DisgoClient api.Disgo
	Listeners   *[]*api.EventListener
	Handlers    *map[string]api.GatewayEventProvider
	Channel     chan api.GenericEvent
}

func (e EventManagerImpl) Disgo() api.Disgo {
	return e.DisgoClient
}

func (e EventManagerImpl) Handle(name string, payload json.RawMessage) {
	if handler, ok := (*e.Handlers)[name]; ok {
		eventPayload := handler.New()
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			log.Errorf("error while unmarshaling event. error: %s", err)
		}
		handler.Handle(e, eventPayload)
	}
}

func (e EventManagerImpl) Dispatch(event api.GenericEvent) {
	e.Channel <- event
}

func (e EventManagerImpl) AddEventListeners(listeners ...api.EventListener) {
	for _, listener := range listeners {
		*e.Listeners = append(*e.Listeners, &listener)
	}
}

func (e EventManagerImpl) ListenEvents() {
	defer func() {
		log.Infof("closing event channel...")
		close(e.Channel)
	}()
	for {
		event := <-e.Channel
		for _, listener := range *e.Listeners {
			(*listener).OnEvent(event)
		}
	}
}

func (e EventManagerImpl) AddHandler(event string, handler api.GatewayEventProvider) {
	(*e.Handlers)[event] = handler
}
