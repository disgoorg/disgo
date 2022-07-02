package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerGuildBanAdd struct{}

func (h *gatewayHandlerGuildBanAdd) EventType() gateway.EventType {
	return gateway.EventTypeGuildBanAdd
}

func (h *gatewayHandlerGuildBanAdd) New() any {
	return &gateway.EventGuildBanAdd{}
}

func (h *gatewayHandlerGuildBanAdd) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventGuildBanAdd)

	client.EventManager().DispatchEvent(&events.GuildBan{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      payload.GuildID,
		User:         payload.User,
	})
}
