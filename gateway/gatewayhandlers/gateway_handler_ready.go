package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerReady handles discord.GatewayEventTypeReady
type gatewayHandlerReady struct{}

// EventType returns the gateway.EventType
func (h *gatewayHandlerReady) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeReady
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerReady) New() interface{} {
	return &discord.GatewayEventReady{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerReady) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	readyEvent := *v.(*discord.GatewayEventReady)

	bot.ApplicationID = readyEvent.Application.ID

	bot.EntityBuilder.CreateSelfUser(readyEvent.SelfUser, core.CacheStrategyYes)

	for _, guild := range readyEvent.Guilds {
		bot.Caches.GuildCache().SetUnready(guild.ID)
	}

	bot.EventManager.Dispatch(&events.ReadyEvent{
		GenericEvent:      events.NewGenericEvent(bot, sequenceNumber),
		GatewayEventReady: readyEvent,
	})

}
