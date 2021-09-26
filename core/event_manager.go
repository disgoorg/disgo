package core

import (
	"io"

	"github.com/DisgoOrg/disgo/discord"
)

// EventManager lets you listen for specific events triggered by raw gateway events
type EventManager interface {
	Bot() *Bot
	Close()
	AddEventListeners(eventListeners ...EventListener)
	RemoveEventListeners(eventListeners ...EventListener)
	HandleGateway(gatewayEventType discord.GatewayEventType, sequenceNumber int, payload io.Reader)
	HandleHTTP(responseChannel chan<- discord.InteractionResponse, payload io.Reader)
	Dispatch(event Event)
}

// EventListener is used to create new EventListener to listen to events
type EventListener interface {
	OnEvent(event interface{})
}

// Event the basic interface each event implement
type Event interface {
	Bot() *Bot
	SequenceNumber() int
}

// GatewayEventHandler is used to handle Gateway Event(s)
type GatewayEventHandler interface {
	EventType() discord.GatewayEventType
	New() interface{}
	HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{})
}

// HTTPServerEventHandler is used to handle HTTP Event(s)
type HTTPServerEventHandler interface {
	New() interface{}
	HandleHTTPEvent(bot *Bot, responseChannel chan<- discord.InteractionResponse, v interface{})
}
