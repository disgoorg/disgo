package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerGuildBanAdd handles discord.GatewayEventTypeIntegrationDelete
type gatewayHandlerIntegrationDelete struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerIntegrationDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeIntegrationDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerIntegrationDelete) New() interface{} {
	return &discord.IntegrationDeleteGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerIntegrationDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.IntegrationDeleteGatewayEvent)

	bot.EventManager.Dispatch(&events.IntegrationDeleteEvent{
		GenericIntegrationEvent: &events.GenericIntegrationEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			GuildId:      payload.GuildID,
		},
		ID:            payload.ID,
		ApplicationID: payload.ApplicationID,
	})
}
