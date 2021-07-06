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
func (h *GuildRoleDeleteHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildRoleDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleDeleteHandler) New() interface{} {
	return &roleCreateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleDeleteHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	roleDeleteData, ok := i.(*roleDeleteData)
	if !ok {
		return
	}

	guild := disgo.Cache().Guild(roleDeleteData.GuildID)
	if guild == nil {
		// todo: replay event later. maybe guild is not cached yet but in a few seconds
		return
	}

	role := disgo.Cache().Role(roleDeleteData.RoleID)
	if role != nil {
		disgo.Cache().UncacheRole(roleDeleteData.GuildID, roleDeleteData.RoleID)
	}

	eventManager.Dispatch(&events.RoleDeleteEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        guild,
			},
			RoleID: roleDeleteData.RoleID,
			Role:   role,
		},
	})
}
