package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

type GenericIntegrationEvent struct {
	*GenericEvent
	GuildId discord.Snowflake
}

// Guild returns the Guild this Integration was created in.
// This will only check cached guilds!
func (e *GenericIntegrationEvent) Guild() *core.Guild {
	return e.Bot().Caches.GuildCache().Get(e.GuildId)
}

// IntegrationCreateEvent indicates that a new Integration was created in a Guild
type IntegrationCreateEvent struct {
	*GenericIntegrationEvent
	Integration *core.Integration
}

// IntegrationUpdateEvent indicates that an integration was updated in a Guild
type IntegrationUpdateEvent struct {
	*GenericIntegrationEvent
	Integration *core.Integration
}

// IntegrationDeleteEvent indicates that an Integration was deleted from a Guild
type IntegrationDeleteEvent struct {
	*GenericIntegrationEvent
	ID            discord.Snowflake
	ApplicationID discord.Snowflake
}

type GuildIntegrationsUpdateEvent struct {
	*GenericIntegrationEvent
}
