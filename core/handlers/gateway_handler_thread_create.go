package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadCreate struct{}

func (h *gatewayHandlerThreadCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadCreate
}

func (h *gatewayHandlerThreadCreate) New() interface{} {
	return &discord.GatewayEventThreadCreate{}
}

func (h *gatewayHandlerThreadCreate) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventThreadCreate)

	bot.EventManager.Dispatch(&events2.ThreadCreateEvent{
		GenericThreadEvent: &events2.GenericThreadEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			ThreadID:     payload.ID(),
			GuildID:      payload.GuildID(),
			Thread:       bot.EntityBuilder.CreateChannel(payload.GuildThread, core.CacheStrategyYes).(core.GuildThread),
		},
	})
}
