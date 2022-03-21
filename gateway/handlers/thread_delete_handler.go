package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

type gatewayHandlerThreadDelete struct{}

func (h *gatewayHandlerThreadDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadDelete
}

func (h *gatewayHandlerThreadDelete) New() any {
	return &discord.GatewayEventThreadDelete{}
}

func (h *gatewayHandlerThreadDelete) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GatewayEventThreadDelete)

	channel, _ := client.Caches().Channels().Remove(payload.ID)
	client.Caches().ThreadMembers().RemoveAll(payload.ID)

	client.EventManager().Dispatch(&events.ThreadDeleteEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber),
			ThreadID:     payload.ID,
			GuildID:      payload.GuildID,
			ParentID:     payload.ParentID,
			Thread:       channel.(discord.GuildThread),
		},
	})
}
