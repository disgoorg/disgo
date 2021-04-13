package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericChannelEvent is called upon receiving any api.Channel api.Event
type GenericChannelEvent struct {
	GenericEvent
	ChannelID api.Snowflake
}

// Channel returns the api.Channel from the api.Cache if cached
func (e GenericChannelEvent) Channel() *api.Channel {
	return e.Disgo().Cache().Channel(e.ChannelID)
}
