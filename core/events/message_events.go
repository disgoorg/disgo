package events

import (
	
	"github.com/DisgoOrg/disgo/discord"
)

// GenericMessageEvent generic api.Message event
type GenericMessageEvent struct {
	*GenericEvent
	MessageID discord.Snowflake
	Message   *core.Message
	ChannelID discord.Snowflake
}

// MessageChannel returns the api.MessageChannel where the GenericMessageEvent happened
func (e *GenericMessageEvent) MessageChannel() *core.MessageChannel {
	return e.Disgo().Cache().MessageChannel(e.ChannelID)
}

// MessageDeleteEvent indicates that an api.Message got deleted
type MessageDeleteEvent struct {
	*GenericMessageEvent
}

// MessageCreateEvent indicates that an api.Message got received
type MessageCreateEvent struct {
	*GenericMessageEvent
}

// MessageUpdateEvent indicates that an api.Message got update
type MessageUpdateEvent struct {
	*GenericMessageEvent
	OldMessage *core.Message
}
