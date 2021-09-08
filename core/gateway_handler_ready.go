package core

import (
	"github.com/DisgoOrg/disgo/discord"
)

// ReadyHandler handles discord.GatewayEventTypeReady
type ReadyHandler struct{}

// EventType returns the gateway.EventType
func (h *ReadyHandler) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeReady
}

// New constructs a new payload receiver for the raw gateway event
func (h *ReadyHandler) New() interface{} {
	return &discord.GatewayEventReady{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ReadyHandler) HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{}) {
	readyEvent := *v.(*discord.GatewayEventReady)

	bot.ApplicationID = readyEvent.Application.ID

	bot.EntityBuilder.CreateSelfUser(readyEvent.SelfUser, CacheStrategyYes)

	for _, guild := range readyEvent.Guilds {
		bot.EntityBuilder.CreateGuild(guild, CacheStrategyYes)
	}

	bot.EventManager.Dispatch(&ReadyEvent{
		GenericEvent:      NewGenericEvent(bot, sequenceNumber),
		GatewayEventReady: readyEvent,
	})

}
