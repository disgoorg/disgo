package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type roleDeleteData struct {
	GuildID discord.Snowflake `json:"guild_id"`
	RoleID  discord.Snowflake `json:"role_id"`
}

// GuildRoleDeleteHandler handles api.GuildRoleDeleteGatewayEvent
type GuildRoleDeleteHandler struct{}

// EventType returns the api.GatewayGatewayEventType
func (h *GuildRoleDeleteHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *GuildRoleDeleteHandler) New() interface{} {
	return roleCreateData{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *GuildRoleDeleteHandler) HandleGatewayEvent(disgo core.Disgo, eventManager core.EventManager, sequenceNumber int, i interface{}) {
	payload, ok := i.(roleDeleteData)
	if !ok {
		return
	}

	role := disgo.Cache().RoleCache().GetCopy(payload.RoleID)

	disgo.Cache().RoleCache().Uncache(payload.GuildID, payload.RoleID)

	eventManager.Dispatch(&events.RoleDeleteEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericGuildEvent: &events.GenericGuildEvent{
				GenericEvent: events.NewGenericEvent(disgo, sequenceNumber),
				Guild:        disgo.Cache().GuildCache().Get(payload.GuildID),
			},
			RoleID: payload.RoleID,
			Role:   role,
		},
	})
}
