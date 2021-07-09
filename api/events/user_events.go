package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericUserEvent is called upon receiving UserUpdateEvent or UserTypingEvent
type GenericUserEvent struct {
	*GenericEvent
	UserID api.Snowflake
	User   *api.User
}

// UserUpdateEvent  indicates that a api.User updated
type UserUpdateEvent struct {
	*GenericUserEvent
	OldUser *api.User
}

// UserTypingEvent indicates that a api.User started typing in a api.DMChannel or api.TextChannel(requires the api.GatewayIntentsDirectMessageTyping and/or api.GatewayIntentsGuildMessageTyping)
type UserTypingEvent struct {
	*GenericUserEvent
	ChannelID api.Snowflake
}

// Channel returns the api.Channel the api.User started typing in
func (e *UserTypingEvent) Channel() api.Channel {
	return e.Disgo().Cache().Channel(e.ChannelID)
}

// MessageChannel returns the api.MessageChannel the api.User started typing in
func (e UserTypingEvent) MessageChannel() api.MessageChannel {
	return e.Disgo().Cache().MessageChannel(e.ChannelID)
}

// DMChannel returns the api.DMChannel the api.User started typing in
func (e *UserTypingEvent) DMChannel() api.DMChannel {
	return e.Disgo().Cache().DMChannel(e.ChannelID)
}

// TextChannel returns the api.TextChannel the api.User started typing in
func (e *UserTypingEvent) TextChannel() api.TextChannel {
	return e.Disgo().Cache().TextChannel(e.ChannelID)
}
