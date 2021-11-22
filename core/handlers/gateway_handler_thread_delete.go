package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	events2 "github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadDelete struct{}

func (h *gatewayHandlerThreadDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadDelete
}

func (h *gatewayHandlerThreadDelete) New() interface{} {
	return &discord.GatewayEventThreadDelete{}
}

func (h *gatewayHandlerThreadDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventThreadDelete)

	channel := bot.Caches.Channels().GetCopy(payload.ID)
	var thread core.GuildThread
	if c, ok := channel.(core.GuildThread); ok {
		thread = c
	}

	bot.EventManager.Dispatch(&events2.ThreadDeleteEvent{
		GenericThreadEvent: &events2.GenericThreadEvent{
			GenericEvent: events2.NewGenericEvent(bot, sequenceNumber),
			ThreadID:     payload.ID,
			GuildID:      payload.GuildID,
			ParentID:     payload.ParentID,
			Thread:       thread,
		},
	})
}
