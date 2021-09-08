package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildInviteEvent is called upon receiving GuildInviteCreateEvent or GuildInviteDeleteEvent(requires core.GatewayIntentsGuildInvites)
type GenericGuildInviteEvent struct {
	*GenericGuildEvent
	Code      string
	ChannelID discord.Snowflake
}

// Channel returns the core.GetChannel the GenericGuildInviteEvent happened in(returns nil if the core.GetChannel is uncached or core.Caches is disabled)
func (e GenericGuildInviteEvent) Channel() *Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}

// GuildInviteCreateEvent is called upon creation of a new core.Invite in an core.Guild(requires core.GatewayIntentsGuildInvites)
type GuildInviteCreateEvent struct {
	*GenericGuildInviteEvent
	Invite *Invite
}

// GuildInviteDeleteEvent is called upon deletion of a new core.Invite in an core.Guild(requires core.GatewayIntentsGuildInvites)
type GuildInviteDeleteEvent struct {
	*GenericGuildInviteEvent
}
