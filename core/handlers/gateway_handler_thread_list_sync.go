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

func (h *gatewayHandlerThreadListSync) New() any {
	return &discord.GatewayEventThreadListSync{}
}

func (h *gatewayHandlerThreadListSync) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GatewayEventThreadListSync)

	for _, thread := range payload.Threads {
		bot.Caches().Channels().Put(thread.ID(), thread)
		bot.EventManager().Dispatch(&events.ThreadShowEvent{
			GenericThreadEvent: &events.GenericThreadEvent{
				Thread:   thread,
				ThreadID: thread.ID(),
				GuildID:  payload.GuildID,
			},
		})
	}
}
