package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildBanAdd handles discord.GatewayEventTypeIntegrationDelete
type gatewayHandlerIntegrationDelete struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerIntegrationDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeIntegrationDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerIntegrationDelete) New() any {
	return &discord.IntegrationDeleteGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerIntegrationDelete) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.IntegrationDeleteGatewayEvent)

	client.EventManager().Dispatch(&events.IntegrationDeleteEvent{
		GenericEvent:  events.NewGenericEvent(client, sequenceNumber),
		GuildID:       payload.GuildID,
		ID:            payload.ID,
		ApplicationID: payload.ApplicationID,
	})
}
