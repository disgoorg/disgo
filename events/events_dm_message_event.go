package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericDMMessageEvent is called upon receiving DMMessageCreateEvent, DMMessageUpdateEvent, DMMessageDeleteEvent, GenericDMMessageReactionEvent, DMMessageReactionAddEvent, DMMessageReactionRemoveEvent, DMMessageReactionRemoveEmojiEvent or DMMessageReactionRemoveAllEvent(requires core.GatewayIntentsDirectMessages)
type GenericDMMessageEvent struct {
	*GenericEvent
	MessageID discord.Snowflake
	Message   *core.Message
	ChannelID discord.Snowflake
}

// Channel returns the core.DMChannel where the GenericDMMessageEvent happened
func (e GenericDMMessageEvent) Channel() *core.DMChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(*core.DMChannel)
	}
	return nil
}

// DMMessageCreateEvent is called upon receiving an core.Message in an core.DMChannel(requires core.GatewayIntentsDirectMessages)
type DMMessageCreateEvent struct {
	*GenericDMMessageEvent
}

// DMMessageUpdateEvent is called upon editing an core.Message in an core.DMChannel(requires core.GatewayIntentsDirectMessages)
type DMMessageUpdateEvent struct {
	*GenericDMMessageEvent
	OldMessage *core.Message
}

// DMMessageDeleteEvent is called upon deleting an core.Message in an core.DMChannel(requires core.GatewayIntentsDirectMessages)
type DMMessageDeleteEvent struct {
	*GenericDMMessageEvent
}
