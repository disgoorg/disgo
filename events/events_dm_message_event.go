package events

import "github.com/DisgoOrg/disgo/core"

// GenericDMMessageEvent is called upon receiving DMMessageCreateEvent, DMMessageUpdateEvent, DMMessageDeleteEvent, GenericDMMessageReactionEvent, DMMessageReactionAddEvent, DMMessageReactionRemoveEvent, DMMessageReactionRemoveEmojiEvent or DMMessageReactionRemoveAllEvent (requires discord.GatewayIntentsDirectMessage)
type GenericDMMessageEvent struct {
	*GenericMessageEvent
}

// DMChannel returns the Channel the GenericDMMessageEvent happened in.
// This will only check cached channels!
func (e GenericDMMessageEvent) DMChannel() *core.Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
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
