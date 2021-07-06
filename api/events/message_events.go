package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericMessageEvent generic api.Message event
type GenericMessageEvent struct {
	*GenericEvent
	MessageID api.Snowflake
	Message   *api.Message
	ChannelID api.Snowflake
}

// MessageChannel returns the api.MessageChannel where the GenericMessageEvent happened
func (e *GenericMessageEvent) MessageChannel() *api.MessageChannel {
	return e.Disgo().Cache().MessageChannel(e.ChannelID)
}

// MessageDeleteEvent indicates that a api.Message got deleted
type MessageDeleteEvent struct {
	*GenericMessageEvent
}

// MessageCreateEvent indicates that a api.Message got received
type MessageCreateEvent struct {
	*GenericMessageEvent
}

// MessageUpdateEvent indicates that a api.Message got update
type MessageUpdateEvent struct {
	*GenericMessageEvent
	OldMessage *api.Message
}
