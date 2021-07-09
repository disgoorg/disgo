package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericRoleEvent generic api.Role event
type GenericRoleEvent struct {
	*GenericGuildEvent
	RoleID api.Snowflake
	Role   *api.Role
}

// RoleCreateEvent indicates that a api.Role got created
type RoleCreateEvent struct {
	*GenericRoleEvent
}

// RoleUpdateEvent indicates that a api.Role got updated
type RoleUpdateEvent struct {
	*GenericRoleEvent
	OldRole *api.Role
}

// RoleDeleteEvent indicates that a api.Role got deleted
type RoleDeleteEvent struct {
	*GenericRoleEvent
}
