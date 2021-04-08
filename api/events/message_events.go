package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericMessageEvent generic api.Message event
type GenericMessageEvent struct {
	GenericEvent
	MessageID api.Snowflake
	ChannelID api.Snowflake
}

// MessageChannel returns the api.MessageChannel where this api.message got received
func (e *GenericMessageEvent) MessageChannel() *api.MessageChannel {
	return e.Disgo().Cache().MessageChannel(e.ChannelID)
}

// MessageDeleteEvent indicates a api.Message got deleted
type MessageDeleteEvent struct {
	GenericMessageEvent
	Message *api.Message
}

// MessageReceivedEvent indicates a api.Message got received
type MessageReceivedEvent struct {
	GenericMessageEvent
	Message *api.Message
}

// MessageUpdateEvent indicates a api.Message got update
type MessageUpdateEvent struct {
	GenericMessageEvent
	NewMessage *api.Message
	OldMessage *api.Message
}
