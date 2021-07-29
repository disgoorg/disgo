package events

import "github.com/DisgoOrg/disgo/api"

// GenericGuildChannelEvent is called upon receiving GuildChannelCreateEvent, GuildChannelUpdateEvent or GuildChannelDeleteEvent
type GenericGuildChannelEvent struct {
	*GenericChannelEvent
	GuildID      api.Snowflake
	GuildChannel *api.GuildChannel
}

// Guild returns the cached api.Guild the event happened in
func (e GenericGuildChannelEvent) Guild() *api.Guild {
	return e.Disgo().Cache().Guild(e.GuildID)
}

// GuildChannelCreateEvent indicates that a new api.GuildChannel got created in an api.Guild
type GuildChannelCreateEvent struct {
	*GenericGuildChannelEvent
}

// GuildChannelUpdateEvent indicates that an api.GuildChannel got updated in an api.Guild
type GuildChannelUpdateEvent struct {
	*GenericGuildChannelEvent
	OldGuildChannel *api.GuildChannel
}

// GuildChannelDeleteEvent indicates that an api.GuildChannel got deleted in an api.Guild
type GuildChannelDeleteEvent struct {
	*GenericGuildChannelEvent
}
