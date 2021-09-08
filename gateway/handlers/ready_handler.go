package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
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
	return discord.GatewayEventReady{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *ReadyHandler) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	readyEvent, ok := v.(discord.GatewayEventReady)
	if !ok {
		return
	}

	bot.ApplicationID = readyEvent.Application.ID

	bot.EntityBuilder.CreateSelfUser(readyEvent.SelfUser, core.CacheStrategyYes)

	for _, guild := range readyEvent.Guilds {
		bot.EntityBuilder.CreateGuild(guild, core.CacheStrategyYes)
	}

	bot.EventManager.Dispatch(&events.ReadyEvent{
		GenericEvent:      events.NewGenericEvent(bot, sequenceNumber),
		GatewayEventReady: readyEvent,
	})

}
