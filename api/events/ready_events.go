package events

import "github.com/DisgoOrg/disgo/api"

// ReadyEvent indicates we received the ReadyEvent from the api.Gateway
type ReadyEvent struct {
	api.GenericEvent
	api.ReadyEventData
}
