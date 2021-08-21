package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
)

type roleDeleteData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	RoleID  discord.Snowflake `json:"role_id"`
}

// GuildRoleDeleteHandler handles api.GuildRoleDeleteGatewayEvent
type GuildRoleDeleteHandler struct{}

// EventType returns the api.GatewayEventType
func (h *GuildRoleDeleteHandler) EventType() gateway.EventType {
	return gateway.EventTypeGuildRoleDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleDeleteHandler) New() interface{} {
	return roleCreateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleDeleteHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	roleDeleteData, ok := i.(roleDeleteData)
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
