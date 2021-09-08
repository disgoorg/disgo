package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// ReadyEvent indicates we received the ReadyEvent from the core.Gateway
type ReadyEvent struct {
	*GenericEvent
	discord.GatewayEventReady
}
