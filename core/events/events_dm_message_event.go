package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// GenericDMMessageEvent is called upon receiving DMMessageCreateEvent, DMMessageUpdateEvent, DMMessageDeleteEvent, GenericDMMessageReactionEvent, DMMessageReactionAddEvent, DMMessageReactionRemoveEvent, DMMessageReactionRemoveEmojiEvent or DMMessageReactionRemoveAllEvent (requires discord.GatewayIntentsDirectMessage)
type GenericDMMessageEvent struct {
	*GenericEvent
	MessageID snowflake.Snowflake
	Message   discord.Message
	ChannelID snowflake.Snowflake
}

// Channel returns the Channel the GenericDMMessageEvent happened in
func (e GenericDMMessageEvent) Channel() (discord.DMChannel, bool) {
	if ch, ok := e.Bot().Caches().Channels().Get(e.ChannelID); ok {
		return ch.(discord.DMChannel), true
	}
	return discord.DMChannel{}, false
}

// DMMessageCreateEvent is called upon receiving a discord.Message in a Channel (requires discord.GatewayIntentsDirectMessage)
type DMMessageCreateEvent struct {
	*GenericDMMessageEvent
}

// DMMessageUpdateEvent is called upon editing a discord.Message in a Channel (requires discord.GatewayIntentsDirectMessage)
type DMMessageUpdateEvent struct {
	*GenericDMMessageEvent
	OldMessage discord.Message
}

// DMMessageDeleteEvent is called upon deleting a discord.Message in a Channel (requires discord.GatewayIntentsDirectMessage)
type DMMessageDeleteEvent struct {
	*GenericDMMessageEvent
}
