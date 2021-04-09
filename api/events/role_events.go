package events

import (
	"github.com/DisgoOrg/disgo/api"
)

// GenericRoleEvent generic api.Role event
type GenericRoleEvent struct {
	GenericGuildEvent
	RoleID api.Snowflake
}

func (e GenericRoleEvent) Role() *api.Role {
	return e.Disgo().Cache().Role(e.RoleID)
}

// RoleCreateEvent indicates that a api.Role got created
type RoleCreateEvent struct {
	GenericGuildEvent
	Role *api.Role
}

// RoleUpdateEvent indicates that a api.Role got updated
type RoleUpdateEvent struct {
	GenericGuildEvent
	NewRole *api.Role
	OldRole *api.Role
}

// RoleDeleteEvent indicates that a api.Role got deleted
type RoleDeleteEvent struct {
	GenericGuildEvent
	Role *api.Role
}
