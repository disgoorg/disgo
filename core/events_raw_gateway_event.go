package core

import (
	"github.com/DisgoOrg/disgo/json"

	"github.com/DisgoOrg/disgo/discord"
)

// RawEvent is called for any core.GatewayGatewayEventType we receive if enabled in the core.BotBuilder/core.Options
type RawEvent struct {
	*GenericEvent
	Type       discord.GatewayEventType
	RawPayload json.RawMessage
}
