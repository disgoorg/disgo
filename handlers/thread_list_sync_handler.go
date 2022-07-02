package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerThreadListSync struct{}

func (h *gatewayHandlerThreadListSync) EventType() gateway.EventType {
	return gateway.EventTypeThreadListSync
}

func (h *gatewayHandlerThreadListSync) New() any {
	return &gateway.EventThreadListSync{}
}

func (h *gatewayHandlerThreadListSync) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventThreadListSync)

	for _, thread := range payload.Threads {
		client.Caches().Channels().Put(thread.ID(), thread)
		client.EventManager().DispatchEvent(&events.ThreadShow{
			GenericThread: &events.GenericThread{
				GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
				Thread:       thread,
				ThreadID:     thread.ID(),
				GuildID:      payload.GuildID,
			},
		})
	}
}
