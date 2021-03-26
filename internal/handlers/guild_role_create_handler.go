package handlers

import (
	"github.com/DiscoOrg/disgo/api"
	"github.com/DiscoOrg/disgo/api/events"
)

// RoleCreateData is the GuildRoleCreate.D payload
type RoleCreateData struct {
	GuildID api.Snowflake `json:"guild_id"`
	Role *api.Role `json:"role"`
}

type RoleCreateHandler struct{}

func (h RoleCreateHandler) New() interface{} {
	return &RoleCreateData{}
}

func (h RoleCreateHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	roleCreateData, ok := i.(*RoleCreateData)
	if !ok {
		return
	}
	roleCreateData.Role.Disgo = disgo
	roleCreateData.Role.GuildID = roleCreateData.GuildID
	disgo.Cache().CacheRole(roleCreateData.Role)

	genericGuildEvent := events.GenericGuildEvent{
		Event:   api.Event{
			Disgo: disgo,
		},
		GuildID: roleCreateData.GuildID,
	}
	eventManager.Dispatch(genericGuildEvent)

	genericRoleEvent := events.GenericGuildRoleEvent{
		GenericGuildEvent: genericGuildEvent,
		Role:              roleCreateData.Role,
		RoleID:            roleCreateData.Role.ID,
	}
	eventManager.Dispatch(genericRoleEvent)

	eventManager.Dispatch(events.GuildRoleCreateEvent{
		GenericGuildEvent: genericGuildEvent,
	})
}