package handlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
)

type gatewayHandlerThreadCreate struct{}

func (h *gatewayHandlerThreadCreate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeThreadCreate
}

func (h *gatewayHandlerThreadCreate) New() any {
	return &discord.GatewayEventThreadCreate{}
}

func (h *gatewayHandlerThreadCreate) HandleGatewayEvent(bot core.Bot, sequenceNumber discord.GatewaySequence, v any) {
	payload := *v.(*discord.GatewayEventThreadCreate)

	bot.Caches().Channels().Put(payload.ID(), payload.GuildThread)
	bot.Caches().ThreadMembers().Put(payload.ID(), payload.ThreadMember.UserID, payload.ThreadMember)

	bot.EventManager().Dispatch(&events.ThreadCreateEvent{
		GenericThreadEvent: &events.GenericThreadEvent{
			GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
			ThreadID:     payload.ID(),
			GuildID:      payload.GuildID(),
			Thread:       payload.GuildThread,
		},
		ThreadMember: payload.ThreadMember,
	})
}
