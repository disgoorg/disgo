package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericDMMessageEvent is called upon receiving DMMessageCreateEvent, DMMessageUpdateEvent, DMMessageDeleteEvent, GenericDMMessageReactionEvent, DMMessageReactionAddEvent, DMMessageReactionRemoveEvent, DMMessageReactionRemoveEmojiEvent or DMMessageReactionRemoveAllEvent(requires api.GatewayIntentsDirectMessages)
type GenericDMMessageEvent struct {
	*GenericMessageEvent
}

// DMChannel returns the api.DMChannel where the GenericDMMessageEvent happened
func (e GenericDMMessageEvent) DMChannel() *api.DMChannel {
	return e.Disgo().Cache().DMChannel(e.ChannelID)
}

// DMMessageCreateEvent is called upon receiving an api.Message in an api.DMChannel(requires api.GatewayIntentsDirectMessages)
type DMMessageCreateEvent struct {
	*GenericDMMessageEvent
}

// DMMessageUpdateEvent is called upon editing an api.Message in an api.DMChannel(requires api.GatewayIntentsDirectMessages)
type DMMessageUpdateEvent struct {
	*GenericDMMessageEvent
	OldMessage *api.Message
}

// DMMessageDeleteEvent is called upon deleting an api.Message in an api.DMChannel(requires api.GatewayIntentsDirectMessages)
type DMMessageDeleteEvent struct {
	*GenericDMMessageEvent
}
