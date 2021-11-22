package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildBanAdd handles discord.GatewayEventTypeIntegrationUpdate
type gatewayHandlerIntegrationUpdate struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerIntegrationUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeIntegrationUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerIntegrationUpdate) New() interface{} {
	return &discord.IntegrationCreateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerIntegrationUpdate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.IntegrationUpdateGatewayEvent)

	bot.EventManager.Dispatch(&events2.IntegrationUpdateEvent{
		GenericIntegrationEvent: &events2.GenericIntegrationEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			GuildID:      payload.GuildID,
			Integration:  bot.EntityBuilder.CreateIntegration(payload.GuildID, payload.Integration, core.CacheStrategyYes),
		},
	})
}
