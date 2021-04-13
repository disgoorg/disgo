package events

import (
	"github.com/DisgoOrg/disgo/api"
)

type GenericUserEvent struct {
	GenericEvent
	UserID api.Snowflake
	User   *api.User
}

type UserUpdateEvent struct {
	GenericUserEvent
	OldUser *api.User
}

type UserTypingEvent struct {
	GenericUserEvent
	ChannelID api.Snowflake
}

func (e UserTypingEvent) Channel() *api.Channel {
	return e.Disgo().Cache().Channel(e.ChannelID)
}
