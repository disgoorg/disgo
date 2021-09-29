package core

import (
	"github.com/DisgoOrg/disgo/json"

	"github.com/DisgoOrg/disgo/discord"
)

// RawEvent is called for any discord.GatewayEventType we receive if enabled in the core.BotBuilder/core.BotConfig
type RawEvent struct {
	*GenericEvent
	Type       discord.GatewayEventType
	RawPayload json.RawMessage
}
