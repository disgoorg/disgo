package events

import "github.com/disgoorg/disgo/gateway"

// Ready indicates we received the Ready from the gateway.Gateway
type Ready struct {
	*GenericEvent
	gateway.EventReady
}

// Resumed indicates disgo resumed the gateway.Gateway
type Resumed struct {
	*GenericEvent
}
