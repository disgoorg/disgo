package events

import (
	"encoding/json"

	"github.com/DisgoOrg/disgo/api"
)

type RawGatewayEvent struct {
	GenericEvent
	Type       api.GatewayEventType
	RawPayload json.RawMessage
	Payload    map[string]interface{}
}
