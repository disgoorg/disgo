package events

import (
	"encoding/json"

	"github.com/DisgoOrg/disgo/discord"
)

// RawEvent is called for any api.GatewayGatewayEventType we receive if enabled in the api.DisgoBuilder/api.Options
type RawEvent struct {
	*GenericEvent
	Type       discord.GatewayEventType
	RawPayload json.RawMessage
}
