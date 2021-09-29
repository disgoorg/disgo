package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildInviteEvent is called upon receiving GuildInviteCreateEvent or GuildInviteDeleteEvent (requires discord.GatewayIntentGuildInvites)
type GenericGuildInviteEvent struct {
	*GenericGuildEvent
	Code      string
	ChannelID discord.Snowflake
}

// Channel returns the Channel the GenericGuildInviteEvent happened in.
// This will only check cached channels!
func (e GenericGuildInviteEvent) Channel() *core.Channel {
	return e.Bot().Caches.ChannelCache().Get(e.ChannelID)
}

// GuildInviteCreateEvent is called upon creation of a new core.Invite in a core.Guild (requires discord.GatewayIntentGuildInvites)
type GuildInviteCreateEvent struct {
	*GenericGuildInviteEvent
	Invite *core.Invite
}

// GuildInviteDeleteEvent is called upon deletion of a core.Invite in a core.Guild (requires discord.GatewayIntentGuildInvites)
type GuildInviteDeleteEvent struct {
	*GenericGuildInviteEvent
}
