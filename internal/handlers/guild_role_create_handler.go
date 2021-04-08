package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type roleCreateData struct {
	GuildID api.Snowflake `json:"guild_id"`
	Role    *api.Role     `json:"role"`
}

// GuildRoleCreateHandler handles api.GuildRoleCreateGatewayEvent
type GuildRoleCreateHandler struct{}

// Event returns the raw gateway event Event
func (h GuildRoleCreateHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildRoleCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildRoleCreateHandler) New() interface{} {
	return &roleCreateData{}
}

// Handle handles the specific raw gateway event
func (h GuildRoleCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	roleCreateData, ok := i.(*roleCreateData)
	if !ok {
		return
	}
	roleCreateData.Role.Disgo = disgo
	roleCreateData.Role.GuildID = roleCreateData.GuildID
	disgo.Cache().CacheRole(roleCreateData.Role)

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		GuildID:      roleCreateData.GuildID,
	}
	eventManager.Dispatch(genericGuildEvent)

	genericRoleEvent := events.GenericRoleEvent{
		GenericGuildEvent: genericGuildEvent,
		RoleID:            roleCreateData.Role.ID,
	}
	eventManager.Dispatch(genericRoleEvent)

	eventManager.Dispatch(events.RoleCreateEvent{
		GenericGuildEvent: genericGuildEvent,
		Role:              roleCreateData.Role,
	})
}
