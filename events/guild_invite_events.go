package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericInvite is called upon receiving InviteCreate or InviteDelete (requires gateway.IntentGuildInvites)
type GenericInvite struct {
	*GenericEvent
	GuildID   *snowflake.ID
	ChannelID snowflake.ID
	Code      string
}

// Channel returns the discord.GuildChannel the GenericInvite happened in.
func (e *GenericInvite) Channel() (discord.GuildChannel, bool) {
	return e.Client().Caches().Channel(e.ChannelID)
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
