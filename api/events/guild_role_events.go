package events

import "github.com/DiscoOrg/disgo/api"

type GenericGuildRoleEvent struct {
	GenericGuildEvent
	Role *api.Role
	RoleID api.Snowflake
}

type GuildRoleCreateEvent struct {
	GenericGuildEvent
}

type GuildRoleDeleteEvent struct {
	GenericGuildEvent
}

type GuildRoleUpdateEvent struct {
	GenericGuildEvent
	OldRole *api.Role
}
