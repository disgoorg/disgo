package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type GenericIntegration struct {
	*GenericEvent
	GuildID     snowflake.ID
	Integration discord.Integration
}

// Guild returns the Guild this Integration was created in.
// This will only check cached guilds!
func (e *GenericIntegration) Guild() (discord.Guild, bool) {
	return e.Client().Caches().Guilds().Get(e.GuildID)
}

// IntegrationCreate indicates that a new Integration was created in a Guild
type IntegrationCreate struct {
	*GenericIntegration
}

// IntegrationUpdate indicates that an integration was updated in a Guild
type IntegrationUpdate struct {
	*GenericIntegration
}

// IntegrationDelete indicates that an Integration was deleted from a Guild
type IntegrationDelete struct {
	*GenericEvent
	ID            snowflake.ID
	GuildID       snowflake.ID
	ApplicationID *snowflake.ID
}

type GuildIntegrationsUpdate struct {
	*GenericEvent
	GuildID snowflake.ID
}

type GuildApplicationCommandPermissionsUpdate struct {
	*GenericEvent
	Permissions discord.ApplicationCommandPermissions
}
