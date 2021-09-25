package core

import "github.com/DisgoOrg/disgo/discord"

type GenericIntegrationEvent struct {
	*GenericEvent
	GuildId discord.Snowflake
}

func (e *GenericIntegrationEvent) Guild() *Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildId)
}

type IntegrationCreateEvent struct {
	*GenericIntegrationEvent
	Integration *Integration
}

type IntegrationUpdateEvent struct {
	*GenericIntegrationEvent
	Integration *Integration
}

type IntegrationDeleteEvent struct {
	*GenericIntegrationEvent
	ID            discord.Snowflake
	ApplicationID discord.Snowflake
}

type GuildIntegrationsUpdateEvent struct {
	*GenericIntegrationEvent
}
