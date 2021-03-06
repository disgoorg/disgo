package disgo

import (
	"reflect"

	"github.com/DiscoOrg/disgo/models/events"
)

type EventHandler interface {
	AddEventHandlers(...func(event events.GenericEvent))
	event(event events.GenericEvent)
}

type eventHandler struct {
	handlers map[reflect.Type][]func(event events.GenericEvent)
	channel  chan events.GenericEvent
}

func (e eventHandler) event(event events.GenericEvent) {
	e.channel <- event
}

func (e eventHandler) AddEventHandlers(handlers ...func(event events.GenericEvent)) {
	for _, handler := range handlers {
		handlerType := reflect.TypeOf(handlers).Field(0).Type
		e.handlers[handlerType] = append(e.handlers[handlerType], handler)
	}
}

func (e eventHandler) Start() {
	go e.listenEvents()
}

func (e eventHandler) listenEvents() {
	defer func() {
		close(e.channel)
	}()
	for {
		event := <-e.channel
		for handlerType, handlers := range e.handlers {
			if handlerType != reflect.TypeOf(event) {
				continue
			}
			for _, handler := range handlers {
				handler(event)
			}
		}
	}
}

func newEventHandler() EventHandler {
	handler := eventHandler{
		channel: make(chan events.GenericEvent),
		handlers: make(map[reflect.Type][]func(event events.GenericEvent), 0),
	}
	handler.Start()
	return handler
}
