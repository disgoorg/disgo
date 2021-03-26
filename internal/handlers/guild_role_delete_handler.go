package handlers

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

// RoleDeleteData is the GuildRoleDelete.D payload
type RoleDeleteData struct {
	GuildID api.Snowflake `json:"guild_id"`
	RoleID  api.Snowflake `json:"role_id"`
}

type RoleDeleteHandler struct{}

func (h RoleDeleteHandler) New() interface{} {
	return &RoleCreateData{}
}

func (h RoleDeleteHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	roleDeleteData, ok := i.(*RoleDeleteData)
	if !ok {
		return
	}

	role := *disgo.Cache().Role(roleDeleteData.GuildID, roleDeleteData.RoleID)
	disgo.Cache().UncacheRole(roleDeleteData.GuildID, roleDeleteData.RoleID)

	genericGuildEvent := events.GenericGuildEvent{
		Event: api.Event{
			Disgo: disgo,
		},
		GuildID: roleDeleteData.GuildID,
	}
	eventManager.Dispatch(genericGuildEvent)

	genericRoleEvent := events.GenericGuildRoleEvent{
		GenericGuildEvent: genericGuildEvent,
		Role:              &role,
		RoleID:            roleDeleteData.RoleID,
	}
	eventManager.Dispatch(genericRoleEvent)

	eventManager.Dispatch(events.GuildRoleDeleteEvent{
		GenericGuildEvent: genericGuildEvent,
	})
}
