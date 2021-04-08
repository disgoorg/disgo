package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericUserEvent struct {
	GenericEvent
	UserID api.Snowflake
}

type UserUpdateEvent struct {
	GenericUserEvent
	NewUser *api.User
	OldUser *api.User
}

type UserTypingEvent struct {
	GenericUserEvent
	User      *api.User
	ChannelID api.Snowflake
}

func (e UserTypingEvent) Channel() *api.Channel {
	return e.Disgo().Cache().Channel(e.ChannelID)
}

type GuildUserTypingEvent struct {
	UserTypingEvent
	GenericGuildMemberEvent
	Member *api.Member
}

func (e UserTypingEvent) TextChannel() *api.TextChannel {
	return e.Disgo().Cache().TextChannel(e.ChannelID)
}

type DMUserTypingEvent struct {
	UserTypingEvent
}

func (e UserTypingEvent) DMChannel() *api.DMChannel {
	return e.Disgo().Cache().DMChannel(e.ChannelID)
}
