package events

import "github.com/DiscoOrg/disgo/api"

// ReadyEvent indicates we received the ReadyEvent from the api.Gateway
type ReadyEvent struct {
	api.Event
	api.ReadyEventData
}
