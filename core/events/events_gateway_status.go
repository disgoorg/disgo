package events

import (
	"github.com/DisgoOrg/disgo/discord"
)

// ReadyEvent indicates we received the ReadyEvent from the core.Gateway
type ReadyEvent struct {
	*GenericEvent
	discord.GatewayEventReady
}

// ResumedEvent indicates disgo resumed to the core.Gateway
type ResumedEvent struct {
	*GenericEvent
}

type InvalidSessionEvent struct {
	*GenericEvent
	MayResume bool
}

// DisconnectedEvent indicates disgo disconnected to the core.Gateway
type DisconnectedEvent struct {
	*GenericEvent
}
