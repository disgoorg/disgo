package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerGuildBanRemove struct{}

func (h *gatewayHandlerGuildBanRemove) EventType() gateway.EventType {
	return gateway.EventTypeGuildBanRemove
}

func (h *gatewayHandlerGuildBanRemove) New() any {
	return &gateway.EventGuildBanRemove{}
}

func (h *gatewayHandlerGuildBanRemove) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventGuildBanRemove)

	client.EventManager().DispatchEvent(&events.GuildUnban{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      payload.GuildID,
		User:         payload.User,
	})
}
