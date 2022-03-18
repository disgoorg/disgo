package events

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/snowflake"
)

// GenericRoleEvent generic core.Role event
type GenericRoleEvent struct {
	*GenericEvent
	GuildID snowflake.Snowflake
	RoleID  snowflake.Snowflake
	Role    discord.Role
}

// RoleCreateEvent indicates that a core.Role got created
type RoleCreateEvent struct {
	*GenericRoleEvent
}

// RoleUpdateEvent indicates that a core.Role got updated
type RoleUpdateEvent struct {
	*GenericRoleEvent
	OldRole discord.Role
}

// RoleDeleteEvent indicates that a core.Role got deleted
type RoleDeleteEvent struct {
	*GenericRoleEvent
}
