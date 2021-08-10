package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/gateway"
)

type roleCreateData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Role    *discord.Role    `json:"role"`
}

// GuildRoleCreateHandler handles api.GuildRoleCreateGatewayEvent
type GuildRoleCreateHandler struct{}

// Event returns the api.GatewayEventType
func (h *GuildRoleCreateHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildRoleCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleCreateHandler) New() interface{} {
	return &roleCreateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleCreateHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	roleCreateData, ok := i.(*roleCreateData)
	if !ok {
		return
	}

	guild := disgo.Cache().Guild(roleCreateData.GuildID)

	role := disgo.EntityBuilder().CreateRole(roleCreateData.GuildID, roleCreateData.Role, core.CacheStrategyYes)

	eventManager.Dispatch(&events.RoleCreateEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        guild,
			},
			RoleID: role.ID,
			Role:   role,
		},
	})
}
