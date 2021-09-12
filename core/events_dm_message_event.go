package core

// GenericDMMessageEvent is called upon receiving DMMessageCreateEvent, DMMessageUpdateEvent, DMMessageDeleteEvent, GenericDMMessageReactionEvent, DMMessageReactionAddEvent, DMMessageReactionRemoveEvent, DMMessageReactionRemoveEmojiEvent or DMMessageReactionRemoveAllEvent(requires core.GatewayIntentsDirectMessages)
type GenericDMMessageEvent struct {
	*GenericMessageEvent
}

// DMChannel returns the core.DMChannel where the GenericDMMessageEvent happened
func (e GenericDMMessageEvent) DMChannel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}

// DMMessageCreateEvent is called upon receiving an core.Message in an core.DMChannel(requires core.GatewayIntentsDirectMessages)
type DMMessageCreateEvent struct {
	*GenericDMMessageEvent
}

// DMMessageUpdateEvent is called upon editing an core.Message in an core.DMChannel(requires core.GatewayIntentsDirectMessages)
type DMMessageUpdateEvent struct {
	*GenericDMMessageEvent
	OldMessage *Message
}

// DMMessageDeleteEvent is called upon deleting an core.Message in an core.DMChannel(requires core.GatewayIntentsDirectMessages)
type DMMessageDeleteEvent struct {
	*GenericDMMessageEvent
}
