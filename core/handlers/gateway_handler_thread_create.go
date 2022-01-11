package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadCreate struct{}

func (h *gatewayHandlerThreadCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadCreate
}

func (h *gatewayHandlerThreadCreate) New() interface{} {
	return &discord.GatewayEventThreadCreate{}
}

func (h *gatewayHandlerThreadCreate) HandleGatewayEvent(bot core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventThreadCreate)

	bot.EventManager().Dispatch(&events.ThreadCreateEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			ThreadID:     payload.ID(),
			GuildID:      payload.GuildID(),
			Thread:       bot.EntityBuilder().CreateChannel(payload.GuildThread, core.CacheStrategyYes).(core.GuildThread),
		},
	})
}
