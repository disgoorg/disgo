package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// gatewayHandlerGuildBanAdd handles discord.GatewayEventTypeIntegrationDelete
type gatewayHandlerIntegrationDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerIntegrationDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeIntegrationDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerIntegrationDelete) New() any {
	return &discord.GatewayEventIntegrationDelete{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerIntegrationDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, v any) {
	payload := *v.(*discord.GatewayEventIntegrationDelete)

	client.EventManager().DispatchEvent(&events.IntegrationDeleteEvent{
		GenericEvent:  events.NewGenericEvent(client, sequenceNumber),
		GuildID:       payload.GuildID,
		ID:            payload.ID,
		ApplicationID: payload.ApplicationID,
	})
}
