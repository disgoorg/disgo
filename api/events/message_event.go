package events

import "github.com/DiscoOrg/disgo/api"

type GenericMessageEvent struct {
	api.Event
	MessageID        api.Snowflake
	MessageChannelID api.Snowflake
}

func (e *GenericMessageEvent) MessageChannel() *api.MessageChannel {
	return e.
		Disgo.
		Cache().
		MessageChannel(e.MessageChannelID)
}


type MessageDeleteEvent struct {
	GenericMessageEvent
	Message api.Message
}


type MessageReceivedEvent struct {
	GenericMessageEvent
	Message api.Message
}


type MessageUpdateEvent struct {
	GenericMessageEvent
	Message api.Message
}

