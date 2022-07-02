package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerThreadDelete struct{}

func (h *gatewayHandlerThreadDelete) EventType() gateway.EventType {
	return gateway.EventTypeThreadDelete
}

func (h *gatewayHandlerThreadDelete) New() any {
	return &gateway.EventThreadDelete{}
}

func (h *gatewayHandlerThreadDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventThreadDelete)

	channel, _ := client.Caches().Channels().Remove(payload.ID)
	client.Caches().ThreadMembers().RemoveAll(payload.ID)

	client.EventManager().DispatchEvent(&events.ThreadDelete{
		GenericThread: &events.GenericThread{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ThreadID:     payload.ID,
			GuildID:      payload.GuildID,
			ParentID:     payload.ParentID,
			Thread:       channel.(discord.GuildThread),
		},
	})
}
