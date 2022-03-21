package handlers

import (
	"github.com/DisgoOrg/disgo/bot"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

type gatewayHandlerThreadCreate struct{}

func (h *gatewayHandlerThreadCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadCreate
}

func (h *gatewayHandlerThreadCreate) New() any {
	return &discord.GatewayEventThreadCreate{}
}

func (h *gatewayHandlerThreadCreate) HandleGatewayEvent(client bot.Client, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GatewayEventThreadCreate)

	client.Caches().Channels().Put(payload.ID(), payload.GuildThread)
	client.Caches().ThreadMembers().Put(payload.ID(), payload.ThreadMember.UserID, payload.ThreadMember)

	client.EventManager().Dispatch(&events.ThreadCreateEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber),
			ThreadID:     payload.ID(),
			GuildID:      payload.GuildID(),
			Thread:       payload.GuildThread,
		},
		ThreadMember: payload.ThreadMember,
	})
}
