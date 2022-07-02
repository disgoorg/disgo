package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericDMMessage is called upon receiving DMMessageCreate , DMMessageUpdate , DMMessageDelete , GenericDMMessageReaction , DMMessageReactionAdd , DMMessageReactionRemove , DMMessageReactionRemoveEmoji or DMMessageReactionRemoveAll (requires gateway.IntentsDirectMessage)
type GenericDMMessage struct {
	*GenericEvent
	MessageID snowflake.ID
	Message   discord.Message
	ChannelID snowflake.ID
}

// Channel returns the Channel the GenericDMMessage happened in
func (e GenericDMMessage) Channel() (discord.DMChannel, bool) {
	if ch, ok := e.Client().Caches().Channels().Get(e.ChannelID); ok {
		return ch.(discord.DMChannel), true
	}
	return discord.DMChannel{}, false
}

// DMMessageCreate is called upon receiving a discord.Message in a Channel (requires gateway.IntentsDirectMessage)
type DMMessageCreate struct {
	*GenericDMMessage
}

// DMMessageUpdate is called upon editing a discord.Message in a Channel (requires gateway.IntentsDirectMessage)
type DMMessageUpdate struct {
	*GenericDMMessage
	OldMessage discord.Message
}

// DMMessageDelete is called upon deleting a discord.Message in a Channel (requires gateway.IntentsDirectMessage)
type DMMessageDelete struct {
	*GenericDMMessage
}
