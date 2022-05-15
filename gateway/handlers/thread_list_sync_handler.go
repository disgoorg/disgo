package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerThreadListSync struct{}

func (h *gatewayHandlerThreadListSync) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadListSync
}

func (h *gatewayHandlerThreadListSync) New() any {
	return &discord.GatewayEventThreadListSync{}
}

func (h *gatewayHandlerThreadListSync) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventThreadListSync)

	for _, thread := range payload.Threads {
		client.Caches().Channels().Put(thread.ID(), thread)
		client.EventManager().DispatchEvent(&events.ThreadShowEvent{
			GenericThreadEvent: &events.GenericThreadEvent{
				Thread:   thread,
				ThreadID: thread.ID(),
				GuildID:  payload.GuildID,
			},
		})
	}
}
