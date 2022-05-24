package events

import (
	"github.com/disgoorg/disgo/discord"
)

// Ready indicates we received the Ready from the gateway.Gateway
type Ready struct {
	*GenericEvent
	discord.GatewayEventReady
}

// Resumed indicates disgo resumed the gateway.Gateway
type Resumed struct {
	*GenericEvent
}
