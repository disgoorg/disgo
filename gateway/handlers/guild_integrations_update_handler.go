package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildEmojisUpdate handles discord.GatewayEventTypeGuildIntegrationsUpdate
type gatewayHandlerGuildIntegrationsUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerGuildIntegrationsUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildIntegrationsUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerGuildIntegrationsUpdate) New() any {
	return &discord.GatewayEventGuildIntegrationsUpdate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerGuildIntegrationsUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildIntegrationsUpdate)

	client.EventManager().DispatchEvent(&events.GuildIntegrationsUpdateEvent{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      payload.GuildID,
	})
}
