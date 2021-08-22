package events

import (
	"github.com/DisgoOrg/disgo/discord"
)

// ReadyEvent indicates we received the ReadyEvent from the api.Gateway
type ReadyEvent struct {
	*GenericEvent
	*discord.ReadyGatewayEvent
}
