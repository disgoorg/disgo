package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMMessageEvent is called upon receiving DMMessageCreateEvent, DMMessageUpdateEvent, DMMessageDeleteEvent, GenericDMMessageReactionEvent, DMMessageReactionAddEvent, DMMessageReactionRemoveEvent, DMMessageReactionRemoveEmojiEvent or DMMessageReactionRemoveAllEvent (requires discord.GatewayIntentsDirectMessage)
type GenericDMMessageEvent struct {
	*GenericEvent
	MessageID discord.Snowflake
	Message   *core.Message
	ChannelID discord.Snowflake
}

// Channel returns the Channel the GenericDMMessageEvent happened in
func (e GenericDMMessageEvent) Channel() *core.DMChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(*core.DMChannel)
	}
	return nil
}

// DMMessageCreateEvent is called upon receiving a core.Message in a Channel (requires discord.GatewayIntentsDirectMessage)
type DMMessageCreateEvent struct {
	*GenericDMMessageEvent
}

// DMMessageUpdateEvent is called upon editing a core.Message in a Channel (requires discord.GatewayIntentsDirectMessage)
type DMMessageUpdateEvent struct {
	*GenericDMMessageEvent
	OldMessage *core.Message
}

// DMMessageDeleteEvent is called upon deleting a core.Message in a Channel (requires discord.GatewayIntentsDirectMessage)
type DMMessageDeleteEvent struct {
	*GenericDMMessageEvent
}
