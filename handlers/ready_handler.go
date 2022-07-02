package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerReady struct{}

func (h *gatewayHandlerReady) EventType() gateway.EventType {
	return gateway.EventTypeReady
}

func (h *gatewayHandlerReady) New() any {
	return &gateway.EventReady{}
}

func (h *gatewayHandlerReady) HandleGatewayEvent(client bot.Client, sequenceNumber int, _ int, v any) {
	readyEvent := *v.(*gateway.EventReady)

	var shardID int
	if readyEvent.Shard != nil {
		shardID = readyEvent.Shard[0]
	}

	client.Caches().PutSelfUser(readyEvent.User)

	for _, guild := range readyEvent.Guilds {
		client.Caches().Guilds().SetUnready(shardID, guild.ID)
	}

	client.EventManager().DispatchEvent(&events.Ready{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		EventReady:   readyEvent,
	})

}
