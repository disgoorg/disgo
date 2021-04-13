package handlers

import (
	"github.com/DisgoOrg/disgo/api"
	"github.com/DisgoOrg/disgo/api/events"
)

type roleUpdateData struct {
	GuildID api.Snowflake `json:"guild_id"`
	Role    *api.Role     `json:"role"`
}

// GuildRoleUpdateHandler handles api.GuildRoleUpdateGatewayEvent
type GuildRoleUpdateHandler struct{}

// Event returns the raw gateway event Event
func (h GuildRoleUpdateHandler) Event() api.GatewayEventType {
	return api.GatewayEventGuildRoleUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h GuildRoleUpdateHandler) New() interface{} {
	return &roleUpdateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h GuildRoleUpdateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	roleUpdateData, ok := i.(*roleUpdateData)
	if !ok {
		return
	}

	guild := disgo.Cache().Guild(roleUpdateData.GuildID)
	if guild == nil {
		// todo: replay event later. maybe guild is not cached yet but in a few seconds
		return
	}

	oldRole := disgo.Cache().Role(roleUpdateData.Role.ID)
	if oldRole != nil {
		oldRole = &*oldRole
	}
	role := disgo.EntityBuilder().CreateRole(roleUpdateData.GuildID, roleUpdateData.Role, api.CacheStrategyYes)

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		Guild:        guild,
	}
	eventManager.Dispatch(genericGuildEvent)

	genericRoleEvent := events.GenericRoleEvent{
		GenericGuildEvent: genericGuildEvent,
		RoleID:            roleUpdateData.Role.ID,
		Role:              role,
	}
	eventManager.Dispatch(genericRoleEvent)

	eventManager.Dispatch(events.RoleUpdateEvent{
		GenericGuildEvent: genericGuildEvent,
		OldRole:           oldRole,
	})
}
