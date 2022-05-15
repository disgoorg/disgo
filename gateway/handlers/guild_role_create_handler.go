package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildRoleCreate handles discord.GatewayEventTypeGuildRoleCreate
type gatewayHandlerGuildRoleCreate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildRoleCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildRoleCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildRoleCreate) New() any {
	return &discord.GatewayEventGuildRoleCreate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildRoleCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildRoleCreate)

	client.Caches().Roles().Put(payload.GuildID, payload.Role.ID, payload.Role)

	client.EventManager().DispatchEvent(&events.RoleCreateEvent{
		GenericRoleEvent: &events.GenericRoleEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			RoleID:       payload.Role.ID,
			Role:         payload.Role,
		},
	})
}
