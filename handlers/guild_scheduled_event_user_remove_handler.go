package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerGuildScheduledEventUserRemove struct{}

func (h *gatewayHandlerGuildScheduledEventUserRemove) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeGuildScheduledEventUserRemove
}

func (h *gatewayHandlerGuildScheduledEventUserRemove) New() any {
	return &discord.GatewayEventGuildScheduledEventUser{}
}

func (h *gatewayHandlerGuildScheduledEventUserRemove) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventGuildScheduledEventUser)

	client.EventManager().DispatchEvent(&events.GuildScheduledEventUserRemove{
		GenericGuildScheduledEventUser: &events.GenericGuildScheduledEventUser{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduledEventID: payload.GuildScheduledEventID,
			UserID:                payload.UserID,
			GuildID:               payload.GuildID,
		},
	})
}
