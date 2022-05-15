package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerThreadDelete struct{}

func (h *gatewayHandlerThreadDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadDelete
}

func (h *gatewayHandlerThreadDelete) New() any {
	return &discord.GatewayEventThreadDelete{}
}

func (h *gatewayHandlerThreadDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventThreadDelete)

	channel, _ := client.Caches().Channels().Remove(payload.ID)
	client.Caches().ThreadMembers().RemoveAll(payload.ID)

	client.EventManager().DispatchEvent(&events.ThreadDeleteEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ThreadID:     payload.ID,
			GuildID:      payload.GuildID,
			ParentID:     payload.ParentID,
			Thread:       channel.(discord.GuildThread),
		},
	})
}
