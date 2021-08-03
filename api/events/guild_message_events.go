package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericGuildMessageEvent is called upon receiving GuildMessageCreateEvent, GuildMessageUpdateEvent or GuildMessageDeleteEvent
type GenericGuildMessageEvent struct {
	*GenericMessageEvent
	GuildID api.Snowflake
}

// Guild returns the api.Guild the GenericGuildMessageEvent happened in
func (e GenericGuildMessageEvent) Guild() *api.Guild {
	return e.Disgo().Cache().Guild(e.GuildID)
}

// TextChannel returns the api.TextChannel from the api.Cache
func (e GenericGuildMessageEvent) TextChannel() *api.TextChannel {
	return e.Disgo().Cache().TextChannel(e.ChannelID)
}

// GuildMessageCreateEvent is called upon receiving an api.Message in an api.DMChannel
type GuildMessageCreateEvent struct {
	*GenericGuildMessageEvent
}

// GuildMessageUpdateEvent is called upon editing an api.Message in an api.DMChannel
type GuildMessageUpdateEvent struct {
	*GenericGuildMessageEvent
	OldMessage *api.Message
}

// GuildMessageDeleteEvent is called upon deleting an api.Message in an api.DMChannel
type GuildMessageDeleteEvent struct {
	*GenericGuildMessageEvent
}
