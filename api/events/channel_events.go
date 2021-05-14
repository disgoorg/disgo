package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericChannelEvent is called upon receiving any api.ChannelImpl api.Event
type GenericChannelEvent struct {
	GenericEvent
	ChannelID api.Snowflake
}

// Channel returns the api.ChannelImpl from the api.Cache if cached
func (e GenericChannelEvent) Channel() api.Channel {
	return e.Disgo().Cache().Channel(e.ChannelID)
}
