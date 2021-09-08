package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericChannelEvent is called upon receiving any api.GetChannel api.EventType
type GenericChannelEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	Channel   *Channel
}
