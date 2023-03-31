package events

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

// GenericRole generic discord.Role event
type GenericRole struct {
	*GenericEvent
	GuildID snowflake.ID
	RoleID  snowflake.ID
	Role    discord.Role
}

// RoleCreate indicates that a discord.Role got created
type RoleCreate struct {
	*GenericRole
}

// RoleUpdate indicates that a discord.Role got updated
type RoleUpdate struct {
	*GenericRole
	OldRole discord.Role
}

// RoleDelete indicates that a discord.Role got deleted
type RoleDelete struct {
	*GenericRole
}
