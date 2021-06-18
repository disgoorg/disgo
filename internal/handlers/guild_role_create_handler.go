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

// HandleGatewayEvent handles the specific raw gateway event
func (h GuildRoleCreateHandler) HandleGatewayEvent(disgo api.Disgo, eventManager api.EventManager, sequenceNumber int, i interface{}) {
	roleCreateData, ok := i.(*roleCreateData)
	if !ok {
		return
	}

	guild := disgo.Cache().Guild(roleCreateData.GuildID)

	role := disgo.EntityBuilder().CreateRole(roleCreateData.GuildID, roleCreateData.Role, api.CacheStrategyYes)

	genericGuildEvent := events.GenericGuildEvent{
		GenericEvent: events.NewEvent(disgo, sequenceNumber),
		Guild:        guild,
	}
	eventManager.Dispatch(genericGuildEvent)

	genericRoleEvent := events.GenericRoleEvent{
		GenericGuildEvent: genericGuildEvent,
		RoleID:            role.ID,
		Role:              role,
	}
	eventManager.Dispatch(genericRoleEvent)

	eventManager.Dispatch(events.RoleCreateEvent{
		GenericGuildEvent: genericGuildEvent,
	})
}
