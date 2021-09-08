package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildInviteEvent is called upon receiving GuildInviteCreateEvent or GuildInviteDeleteEvent(requires api.GatewayIntentsGuildInvites)
type GenericGuildInviteEvent struct {
	*GenericGuildEvent
	Code      string
	ChannelID discord.Snowflake
}

// Channel returns the api.GetChannel the GenericGuildInviteEvent happened in(returns nil if the api.GetChannel is uncached or api.Caches is disabled)
func (e GenericGuildInviteEvent) Channel() *core.Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}

// GuildInviteCreateEvent is called upon creation of a new api.Invite in an api.Guild(requires api.GatewayIntentsGuildInvites)
type GuildInviteCreateEvent struct {
	*GenericGuildInviteEvent
	Invite *core.Invite
}

// GuildInviteDeleteEvent is called upon deletion of a new api.Invite in an api.Guild(requires api.GatewayIntentsGuildInvites)
type GuildInviteDeleteEvent struct {
	*GenericGuildInviteEvent
}
