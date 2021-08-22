package events

import (
	"encoding/json"
)

// RawEvent is called for any api.GatewayGatewayEventType we receive if enabled in the api.DisgoBuilder/api.Options
type RawEvent struct {
	*GenericEvent
	Type       string
	RawPayload json.RawMessage
	Payload    map[string]interface{}
}
