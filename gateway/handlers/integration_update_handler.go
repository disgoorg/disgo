package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildBanAdd handles discord.GatewayEventTypeIntegrationUpdate
type gatewayHandlerIntegrationUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerIntegrationUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeIntegrationUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerIntegrationUpdate) New() any {
	return &discord.IntegrationCreateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerIntegrationUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.IntegrationUpdateGatewayEvent)

	client.EventManager().Dispatch(&events.IntegrationUpdateEvent{
		GenericIntegrationEvent: &events.GenericIntegrationEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber),
			GuildID:      payload.GuildID,
			Integration:  payload.Integration,
		},
	})
}
