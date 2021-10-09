package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
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

	bot.EventManager.Dispatch(&events.IntegrationUpdateEvent{
		GenericIntegrationEvent: &events.GenericIntegrationEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			GuildId:      payload.GuildID,
		},
		Integration: bot.EntityBuilder.CreateIntegration(payload.GuildID, payload.Integration, core.CacheStrategyYes),
	})
}
