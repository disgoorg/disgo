package events

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
)

// GenericRoleEvent generic core.Role event
type GenericRoleEvent struct {
	*GenericEvent
	GuildID discord.Snowflake
	RoleID  discord.Snowflake
	Role    *core.Role
}

// RoleCreateEvent indicates that a core.Role got created
type RoleCreateEvent struct {
	*GenericRoleEvent
}

// RoleUpdateEvent indicates that a core.Role got updated
type RoleUpdateEvent struct {
	*GenericRoleEvent
	OldRole *core.Role
}

// RoleDeleteEvent indicates that a core.Role got deleted
type RoleDeleteEvent struct {
	*GenericRoleEvent
}
