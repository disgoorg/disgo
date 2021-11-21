package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildInviteEvent is called upon receiving GuildInviteCreateEvent or GuildInviteDeleteEvent(requires core.GatewayIntentsGuildInvites)
type GenericGuildInviteEvent struct {
	*GenericEvent
	GuildID   discord.Snowflake
	ChannelID discord.Snowflake
	Code      string
}

// Channel returns the core.GuildChannel the GenericGuildInviteEvent happened in(returns nil if the core.GetChannel is uncached or core.Caches is disabled)
func (e GenericGuildInviteEvent) Channel() core.GuildChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(core.GuildChannel)
	}
	return nil
}

// GuildInviteCreateEvent is called upon creation of a new core.Invite in an core.Guild(requires core.GatewayIntentsGuildInvites)
type GuildInviteCreateEvent struct {
	*GenericGuildInviteEvent
	Invite *core.Invite
}

// GuildInviteDeleteEvent is called upon deletion of a new core.Invite in an core.Guild(requires core.GatewayIntentsGuildInvites)
type GuildInviteDeleteEvent struct {
	*GenericGuildInviteEvent
}
