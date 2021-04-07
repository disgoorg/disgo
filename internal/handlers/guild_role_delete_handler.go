package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type roleDeleteData struct {
	GuildID api.Snowflake `json:"guild_id"`
	RoleID  api.Snowflake `json:"role_id"`
}

// GuildRoleDeleteHandler handles api.GuildRoleDeleteGatewayEvent
type GuildRoleDeleteHandler struct{}

// Event returns the raw gateway event Event
func (h GuildRoleDeleteHandler) Event() api.GatewayEventName {
	return api.GatewayEventGuildRoleDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildRoleDeleteHandler) New() interface{} {
	return &roleCreateData{}
}

// Handle handles the specific raw gateway event
func (h GuildRoleDeleteHandler) Handle(disgo api.Disgo, eventManager api.EventManager, i interface{}) {
	roleDeleteData, ok := i.(*roleDeleteData)
	if !ok {
		return
	}

	role := *disgo.Cache().Role(roleDeleteData.GuildID, roleDeleteData.RoleID)
	disgo.Cache().UncacheRole(roleDeleteData.GuildID, roleDeleteData.RoleID)

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: events.NewEvent(disgo),
		GuildID:      roleDeleteData.GuildID,
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
