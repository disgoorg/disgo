package events

import (
	"encoding/json"

	"github.com/DisgoOrg/disgo/api"
)

// RawGatewayEvent is called for any api.GatewayEventType we receive if enabled in the api.DisgoBuilder/api.Options
type RawGatewayEvent struct {
	*GenericEvent
	Type       api.GatewayEventType
	RawPayload json.RawMessage
	Payload    map[string]interface{}
}
