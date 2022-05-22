package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildBanAdd handles discord.GatewayEventTypeIntegrationCreate
type gatewayHandlerIntegrationCreate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerIntegrationCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeIntegrationCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerIntegrationCreate) New() any {
	return &discord.GatewayEventIntegrationCreate{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerIntegrationCreate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventIntegrationCreate)

	client.EventManager().DispatchEvent(&events.IntegrationCreate{
		GenericIntegration: &events.GenericIntegration{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			Integration:  payload.Integration,
		},
	})
}
