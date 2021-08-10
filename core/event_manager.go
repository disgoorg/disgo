package core

import (
	"io"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/httpserver"
)

// EventManager lets you listen for specific events triggered by raw gateway events
type EventManager interface {
	Disgo() Disgo
	Close()
	AddEventListeners(eventListeners ...EventListener)
	RemoveEventListener(eventListeners ...EventListener)
	HandleGateway(eventType gateway.EventType, sequenceNumber int, payload io.Reader)
	HandleHTTP(eventType httpserver.EventType, replyChannel chan discord.InteractionResponse, payload io.Reader)
	Dispatch(event Event)
}

// EventListener is used to create new EventListener to listen to events
type EventListener interface {
	OnEvent(event interface{})
}

// Event the basic interface each event implement
type Event interface {
	Disgo() Disgo
	SequenceNumber() int
}

// GatewayEventHandler is used to handle Gateway Event(s)
type GatewayEventHandler interface {
	EventType() gateway.EventType
	New() interface{}
	HandleGatewayEvent(disgo Disgo, eventManager EventManager, sequenceNumber int, payload interface{})
}

// HTTPEventHandler is used to handle HTTP EventType(s)
type HTTPEventHandler interface {
	EventType() httpserver.EventType
	New() interface{}
	HandleHTTPEvent(disgo Disgo, eventManager EventManager, replyChannel chan discord.InteractionResponse, payload interface{})
}
