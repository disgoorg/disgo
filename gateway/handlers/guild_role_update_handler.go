package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type roleUpdateData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Role    discord.Role      `json:"role"`
}

// GuildRoleUpdateHandler handles api.GuildRoleUpdateGatewayEvent
type GuildRoleUpdateHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildRoleUpdateHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleUpdateHandler) New() interface{} {
	return &roleUpdateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleUpdateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, v interface{}) {
	payload, ok := v.(roleUpdateData)
	if !ok {
		return
	}

	oldRole := disgo.Caches().RoleCache().GetCopy(payload.Role.ID)

	eventManager.Dispatch(&events.RoleUpdateEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        disgo.Caches().GuildCache().Get(payload.GuildID),
			},
			RoleID: payload.Role.ID,
			Role:   disgo.EntityBuilder().CreateRole(payload.GuildID, payload.Role, core.CacheStrategyYes),
		},
		OldRole: oldRole,
	})
}
