package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadListSync struct{}

func (h *gatewayHandlerThreadListSync) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadListSync
}

func (h *gatewayHandlerThreadListSync) New() interface{} {
	return &discord.GatewayEventThreadListSync{}
}

func (h *gatewayHandlerThreadListSync) HandleGatewayEvent(bot core.Bot, sequenceNumber int, v interface{}) {
	payload := *v.(*discord.GatewayEventThreadListSync)

	for i := range payload.Threads {
		thread := bot.EntityBuilder().CreateChannel(payload.Threads[i], core.CacheStrategyYes).(core.GuildThread)
		bot.EventManager().Dispatch(&events.ThreadShowEvent{
			GenericThreadEvent: &events.GenericThreadEvent{
				Thread:   thread,
				ThreadID: thread.ID(),
				GuildID:  payload.GuildID,
			},
		})
	}
}
