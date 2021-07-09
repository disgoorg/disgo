package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericChannelEvent is called upon receiving any api.Channel api.Event
type GenericChannelEvent struct {
	*GenericEvent
	ChannelID api.Snowflake
	Channel   api.Channel
}
