package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

// InviteCreate is called upon creation of a new discord.Invite (requires gateway.IntentGuildInvites)
type InviteCreate struct {
	*GenericEvent

	gateway.EventInviteCreate
}

// Channel returns the discord.GuildChannel the GenericInvite happened in.
func (e *InviteCreate) Channel() (discord.GuildChannel, bool) {
	return e.Client().Caches().Channel(e.ChannelID)
}

func (e *InviteCreate) Guild() (discord.Guild, bool) {
	if e.GuildID == nil {
		return discord.Guild{}, false
	}
	return e.Client().Caches().Guild(*e.GuildID)
}

// InviteDelete is called upon deletion of a discord.Invite (requires gateway.IntentGuildInvites)
type InviteDelete struct {
	*GenericEvent

	GuildID   *snowflake.ID
	ChannelID snowflake.ID
	Code      string
}

// Channel returns the discord.GuildChannel the GenericInvite happened in.
func (e *InviteDelete) Channel() (discord.GuildChannel, bool) {
	return e.Client().Caches().Channel(e.ChannelID)
}

func (e *InviteDelete) Guild() (discord.Guild, bool) {
	if e.GuildID == nil {
		return discord.Guild{}, false
	}
	return e.Client().Caches().Guild(*e.GuildID)
}
