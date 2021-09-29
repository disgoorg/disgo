package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// gatewayHandlerGuildBanAdd handles discord.GatewayEventTypeIntegrationUpdate
type gatewayHandlerIntegrationUpdate struct{}

// EventType returns the discord.GatewayEventType
func (h *gatewayHandlerIntegrationUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeIntegrationUpdate
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerIntegrationUpdate) New() interface{} {
	return &discord.IntegrationCreateGatewayEvent{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerIntegrationUpdate) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.IntegrationUpdateGatewayEvent)

	bot.EventManager.Dispatch(&IntegrationUpdateEvent{
		GenericIntegrationEvent: &GenericIntegrationEvent{
			GenericEvent: NewGenericEvent(bot, sequenceNumber),
			GuildId:      payload.GuildID,
		},
		Integration: bot.EntityBuilder.CreateIntegration(payload.GuildID, payload.Integration, CacheStrategyYes),
	})
}
