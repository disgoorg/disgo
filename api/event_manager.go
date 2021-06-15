package api

import "encoding/json"

// EventManager lets you listen for specific events triggered by raw gateway events
type EventManager interface {
	Disgo() Disgo
	Close()
	AddEventListeners(eventListeners ...EventListener)
	Handle(eventType GatewayEventType, replyChannel chan InteractionResponse, sequenceNumber int, payload json.RawMessage)
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

// EventHandler provides info about the EventHandler
type EventHandler interface {
	Event() GatewayEventType
	New() interface{}
}

// GatewayEventHandler is used to handle raw gateway events
type GatewayEventHandler interface {
	EventHandler
	HandleGatewayEvent(disgo Disgo, eventManager EventManager, sequenceNumber int, payload interface{})
}

// WebhookEventHandler is used to handle raw webhook events
type WebhookEventHandler interface {
	EventHandler
	HandleWebhookEvent(disgo Disgo, eventManager EventManager, replyChannel chan InteractionResponse, payload interface{})
}
