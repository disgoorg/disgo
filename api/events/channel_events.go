package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericChannelEvent is called upon receiving an event in a api.Channel
type GenericChannelEvent struct {
	GenericEvent
	ChannelID api.Snowflake
}

func (e GenericChannelEvent) Channel() *api.Channel {
	return e.Disgo().Cache().Channel(e.ChannelID)
}
