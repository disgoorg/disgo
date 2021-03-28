package events

import "github.com/DisgoOrg/disgo/api"

// GenericGuildRoleEvent generic api.Role event
type GenericGuildRoleEvent struct {
	GenericGuildEvent
	Role   *api.Role
	RoleID api.Snowflake
}

// GuildRoleCreateEvent indicates that a api.Role got created
type GuildRoleCreateEvent struct {
	GenericGuildEvent
}

// GuildRoleDeleteEvent indicates that a api.Role got deleted
type GuildRoleDeleteEvent struct {
	GenericGuildEvent
}

// GuildRoleUpdateEvent indicates that a api.Role got updated
type GuildRoleUpdateEvent struct {
	GenericGuildEvent
	OldRole *api.Role
}
