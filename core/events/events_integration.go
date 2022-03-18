package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

type GenericIntegrationEvent struct {
	*GenericEvent
	GuildID     snowflake.Snowflake
	Integration discord.Integration
}

// Guild returns the Guild this Integration was created in.
// This will only check cached guilds!
func (e *GenericIntegrationEvent) Guild() (discord.Guild, bool) {
	return e.Bot().Caches().Guilds().Get(e.GuildID)
}

// IntegrationCreateEvent indicates that a new Integration was created in a Guild
type IntegrationCreateEvent struct {
	*GenericIntegrationEvent
}

// IntegrationUpdateEvent indicates that an integration was updated in a Guild
type IntegrationUpdateEvent struct {
	*GenericIntegrationEvent
}

// IntegrationDeleteEvent indicates that an Integration was deleted from a Guild
type IntegrationDeleteEvent struct {
	*GenericEvent
	ID            snowflake.Snowflake
	GuildID       snowflake.Snowflake
	ApplicationID *snowflake.Snowflake
}

type GuildIntegrationsUpdateEvent struct {
	*GenericEvent
	GuildID snowflake.Snowflake
}
