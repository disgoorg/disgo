package events

import (
	
	"github.com/DisgoOrg/disgo/discord"
)

// GenericChannelEvent is called upon receiving any api.Channel api.EventType
type GenericChannelEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	Channel   *core.Channel
}
