package events

import (
	"github.com/DisgoOrg/disgo/discord"
)

// ReadyEvent indicates we received the ReadyEvent from the gateway.Gateway
type ReadyEvent struct {
	*GenericEvent
	discord.GatewayEventReady
}

// ResumedEvent indicates disgo resumed the gateway.Gateway
type ResumedEvent struct {
	*GenericEvent
}

type InvalidSessionEvent struct {
	*GenericEvent
	MayResume bool
}

// DisconnectedEvent indicates disgo disconnected from the gateway.Gateway
type DisconnectedEvent struct {
	*GenericEvent
}
