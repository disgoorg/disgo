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

// Guild returns the Guild this Integration was created in.
// This will only check cached guilds!
func (e *GenericIntegrationEvent) Guild() *core.Guild {
	return e.Bot().Caches.Guilds().Get(e.GuildID)
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
	ID            discord.Snowflake
	GuildID       discord.Snowflake
	ApplicationID *discord.Snowflake
}

type GuildIntegrationsUpdateEvent struct {
	*GenericEvent
	GuildID discord.Snowflake
}
