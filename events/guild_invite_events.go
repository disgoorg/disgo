package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericInvite is called upon receiving InviteCreate or InviteDelete (requires gateway.IntentGuildInvites)
type GenericInvite struct {
	*GenericEvent
	GuildID   *snowflake.ID
	ChannelID snowflake.ID
	Code      string
}

// Channel returns the Channel the GenericInvite happened in.
func (e GenericInvite) Channel() (discord.GuildChannel, bool) {
	return e.Client().Caches().Channels().GetGuildChannel(e.ChannelID)
}

// InviteCreate is called upon creation of a new discord.Invite (requires gateway.IntentGuildInvites)
type InviteCreate struct {
	*GenericInvite
	Invite discord.Invite
}

// InviteDelete is called upon deletion of a discord.Invite (requires gateway.IntentGuildInvites)
type InviteDelete struct {
	*GenericInvite
}
