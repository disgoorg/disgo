package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericGuildInviteEvent is called upon receiving GuildInviteCreateEvent or GuildInviteDeleteEvent (requires discord.GatewayIntentGuildInvites)
type GenericGuildInviteEvent struct {
	*GenericEvent
	GuildID   discord.Snowflake
	ChannelID discord.Snowflake
	Code      string
}

// Channel returns the Channel the GenericGuildInviteEvent happened in.
func (e GenericGuildInviteEvent) Channel() core.GuildChannel {
	if ch := e.Bot().Caches.Channels().Get(e.ChannelID); ch != nil {
		return ch.(core.GuildChannel)
	}
	return nil
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
