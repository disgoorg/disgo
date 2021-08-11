package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

type roleUpdateData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Role    *discord.Role     `json:"role"`
}

// GuildRoleUpdateHandler handles api.GuildRoleUpdateGatewayEvent
type GuildRoleUpdateHandler struct{}

// Event returns the api.GatewayEventType
func (h *GuildRoleUpdateHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildRoleUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleUpdateHandler) New() interface{} {
	return &roleUpdateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
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
	role := disgo.EntityBuilder().CreateRole(roleUpdateData.GuildID, roleUpdateData.Role, core.CacheStrategyYes)

	eventManager.Dispatch(&events.RoleUpdateEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        guild,
			},
			RoleID: roleUpdateData.Role.ID,
			Role:   role,
		},
		OldRole: oldRole,
	})
}
