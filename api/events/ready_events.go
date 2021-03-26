package events

import "github.com/DiscoOrg/disgo/api"

type ReadyEvent struct {
	api.Event
	api.ReadyEventData
}
