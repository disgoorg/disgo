package events

import (
	"github.com/DisgoOrg/snowflake"
	"github.com/disgoorg/disgo/discord"
)

// GenericInviteEvent is called upon receiving InviteCreateEvent or InviteDeleteEvent (requires discord.GatewayIntentGuildInvites)
type GenericInviteEvent struct {
	*GenericEvent
	GuildID   *snowflake.Snowflake
	ChannelID snowflake.Snowflake
	Code      string
}

// Channel returns the Channel the GenericInviteEvent happened in.
func (e GenericInviteEvent) Channel() (discord.GuildChannel, bool) {
	return e.Client().Caches().Channels().GetGuildChannel(e.ChannelID)
}

// InviteCreateEvent is called upon creation of a new discord.Invite (requires discord.GatewayIntentGuildInvites)
type InviteCreateEvent struct {
	*GenericInviteEvent
	Invite discord.Invite
}

// InviteDeleteEvent is called upon deletion of a discord.Invite (requires discord.GatewayIntentGuildInvites)
type InviteDeleteEvent struct {
	*GenericInviteEvent
}
