package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// GenericRoleEvent generic core.Role event
type GenericRoleEvent struct {
	*GenericGuildEvent
	RoleID discord.Snowflake
	Role   *Role
}

// RoleCreateEvent indicates that a core.Role got created
type RoleCreateEvent struct {
	*GenericRoleEvent
}

// RoleUpdateEvent indicates that a core.Role got updated
type RoleUpdateEvent struct {
	*GenericRoleEvent
	OldRole *Role
}

// RoleDeleteEvent indicates that a core.Role got deleted
type RoleDeleteEvent struct {
	*GenericRoleEvent
}
