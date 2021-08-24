package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericChannelEvent is called upon receiving any api.GetChannel api.EventType
type GenericChannelEvent struct {
	*GenericEvent
	ChannelID discord.Snowflake
	Channel   core.Channel
}

type ChannelPinsUpdateEvent struct {
	*GenericChannelEvent
	GuildID             *discord.Snowflake
	OldLastPinTimestamp *discord.Time
	NewLastPinTimestamp *discord.Time
}
