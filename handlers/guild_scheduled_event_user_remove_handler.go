package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildScheduledEventUserRemove struct {}

func (h *gatewayHandlerGuildScheduledEventUserRemove) EventType() gateway.EventType {
	return gateway.EventTypeGuildScheduledEventUserRemove
}

func (h *gatewayHandlerGuildScheduledEventUserRemove) New() any {
	return &gateway.EventGuildScheduledEventUserRemove{}
}

func (h *gatewayHandlerGuildScheduledEventUserRemove) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventGuildScheduledEventUserRemove)

	client.EventManager().DispatchEvent(&events.GuildScheduledEventUserRemove{
		GenericGuildScheduledEventUser: &events.GenericGuildScheduledEventUser{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduledEventID: payload.GuildScheduledEventID,
			UserID:                payload.UserID,
			GuildID:               payload.GuildID,
		},
	})
}
