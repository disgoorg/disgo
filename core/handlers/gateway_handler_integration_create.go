package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildBanAdd handles discord.GatewayEventTypeIntegrationCreate
type gatewayHandlerIntegrationCreate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerIntegrationCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeIntegrationCreate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerIntegrationCreate) New() interface{} {
	return &discord.IntegrationCreateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerIntegrationCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.IntegrationCreateGatewayEvent)

	bot.EventManager.Dispatch(&events2.IntegrationCreateEvent{
		GenericIntegrationEvent: &events2.GenericIntegrationEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			Integration:  bot.EntityBuilder.CreateIntegration(payload.GuildID, payload.Integration, core.CacheStrategyYes),
		},
	})
}
