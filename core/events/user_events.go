package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericUserEvent is called upon receiving UserUpdateEvent or UserTypingEvent
type GenericUserEvent struct {
	*GenericEvent
	UserID discord.Snowflake
	User   *core.User
}

// UserUpdateEvent  indicates that an api.User updated
type UserUpdateEvent struct {
	*GenericUserEvent
	OldUser *core.User
}

// UserTypingEvent indicates that an api.User started typing in an api.DMChannel or api.TextChannel(requires the api.GatewayIntentsDirectMessageTyping and/or api.GatewayIntentsGuildMessageTyping)
type UserTypingEvent struct {
	*GenericUserEvent
	ChannelID discord.Snowflake
}

// Channel returns the api.GetChannel the api.User started typing in
func (e *UserTypingEvent) Channel() core.Channel {
	return e.Disgo().Caches().ChannelCache().GetChannel(e.ChannelID)
}

// DMChannel returns the api.DMChannel the api.User started typing in
func (e *UserTypingEvent) DMChannel() core.DMChannel {
	return e.Disgo().Caches().DMChannelCache().Get(e.ChannelID)
}

// TextChannel returns the api.TextChannel the api.User started typing in
func (e *UserTypingEvent) TextChannel() core.TextChannel {
	return e.Disgo().Caches().TextChannelCache().Get(e.ChannelID)
}
