package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
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

	bot.EventManager.Dispatch(&events.ThreadDeleteEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			ThreadID:     payload.ID,
			GuildID:      payload.GuildID,
			ParentID:     payload.ParentID,
			Thread:       thread,
		},
	})
}
