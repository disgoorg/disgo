package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// GenericRoleEvent generic discord.Role event
type GenericRoleEvent struct {
	*GenericEvent
	GuildID snowflake.ID
	RoleID  snowflake.ID
	Role    discord.Role
}

// RoleCreateEvent indicates that a discord.Role got created
type RoleCreateEvent struct {
	*GenericRoleEvent
}

// RoleUpdateEvent indicates that a discord.Role got updated
type RoleUpdateEvent struct {
	*GenericRoleEvent
	OldRole discord.Role
}

// RoleDeleteEvent indicates that a discord.Role got deleted
type RoleDeleteEvent struct {
	*GenericRoleEvent
}
