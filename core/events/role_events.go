package events

import (
	
	"github.com/DisgoOrg/disgo/discord"
)

// GenericRoleEvent generic api.Role event
type GenericRoleEvent struct {
	*GenericGuildEvent
	RoleID discord.Snowflake
	Role   *core.Role
}

// RoleCreateEvent indicates that an api.Role got created
type RoleCreateEvent struct {
	*GenericRoleEvent
}

// RoleUpdateEvent indicates that an api.Role got updated
type RoleUpdateEvent struct {
	*GenericRoleEvent
	OldRole *core.Role
}

// RoleDeleteEvent indicates that an api.Role got deleted
type RoleDeleteEvent struct {
	*GenericRoleEvent
}
