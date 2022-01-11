package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildRoleCreate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildRoleCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildRoleCreate) New() interface{} {
	return &discord.GuildRoleCreateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildRoleCreate) HandleGatewayEvent(bot core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GuildRoleCreateGatewayEvent)

	bot.EventManager().Dispatch(&events.RoleCreateEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			RoleID:       payload.Role.ID,
			Role:         bot.EntityBuilder().CreateRole(payload.GuildID, payload.Role, core.CacheStrategyYes),
		},
	})
}
