package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericGuildMessageEvent indicates that we received a api.Message api.GenericEvent in a api.Guild
type GenericGuildMessageEvent struct {
	GenericMessageEvent
	GuildID api.Snowflake
}

// Guild returns the api.Guild from the api.Cache
func (e GenericGuildMessageEvent) Guild() *api.Guild {
	return e.Disgo().Cache().Guild(e.GuildID)
}

// GuildMessageReceivedEvent indicates that we received a api.Message in a api.Guild
type GuildMessageReceivedEvent struct {
	GenericGuildMessageEvent
	Message api.Message
}

// GuildMessageUpdateEvent indicates that a api.Message was updated in a api.Guild
type GuildMessageUpdateEvent struct {
	GenericGuildMessageEvent
	Message api.Message
}

// GuildMessageDeleteEvent indicates that a api.Message was deleted in a api.Guild
type GuildMessageDeleteEvent struct {
	GenericGuildMessageEvent
	Message *api.Message
}
