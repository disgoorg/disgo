package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// GenericGuildInviteEvent is called upon receiving GuildInviteCreateEvent or GuildInviteDeleteEvent (requires discord.GatewayIntentGuildInvites)
type GenericGuildInviteEvent struct {
	*GenericEvent
	GuildID   snowflake.Snowflake
	ChannelID snowflake.Snowflake
	Code      string
}

// Channel returns the Channel the GenericGuildInviteEvent happened in.
func (e GenericGuildInviteEvent) Channel() (discord.GuildChannel, bool) {
	return e.Bot().Caches().Channels().GetGuildChannel(e.ChannelID)
}

// GuildInviteCreateEvent is called upon creation of a new discord.Invite in a discord.Guild (requires discord.GatewayIntentGuildInvites)
type GuildInviteCreateEvent struct {
	*GenericGuildInviteEvent
	Invite discord.Invite
}

// GuildInviteDeleteEvent is called upon deletion of a discord.Invite in a discord.Guild (requires discord.GatewayIntentGuildInvites)
type GuildInviteDeleteEvent struct {
	*GenericGuildInviteEvent
}
