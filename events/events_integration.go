package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type GenericIntegrationEvent struct {
	*GenericEvent
	GuildID     discord.Snowflake
	Integration core.Integration
}

func (e *GenericIntegrationEvent) Guild() *core.Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildID)
}

type IntegrationCreateEvent struct {
	*GenericIntegrationEvent
}

type IntegrationUpdateEvent struct {
	*GenericIntegrationEvent
}

type IntegrationDeleteEvent struct {
	*GenericEvent
	ID            discord.Snowflake
	GuildID       discord.Snowflake
	ApplicationID *discord.Snowflake
}

type GuildIntegrationsUpdateEvent struct {
	*GenericEvent
	GuildID discord.Snowflake
}
