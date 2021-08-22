package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

type roleCreateData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Role    discord.Role      `json:"role"`
}

// GuildRoleCreateHandler handles api.GuildRoleCreateGatewayEvent
type GuildRoleCreateHandler struct{}

// EventType returns the api.GatewayEventType
func (h *GuildRoleCreateHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildRoleCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleCreateHandler) New() interface{} {
	return &roleCreateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleCreateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(roleCreateData)
	if !ok {
		return
	}

	eventManager.Dispatch(&events.RoleCreateEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        disgo.Cache().GuildCache().Get(payload.GuildID),
			},
			RoleID: payload.Role.ID,
			Role:   disgo.EntityBuilder().CreateRole(payload.GuildID, payload.Role, core.CacheStrategyYes),
		},
	})
}
