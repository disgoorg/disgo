package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerWebhooksUpdate struct{}

func (h *gatewayHandlerWebhooksUpdate) EventType() gateway.EventType {
	return gateway.EventTypeWebhooksUpdate
}

func (h *gatewayHandlerWebhooksUpdate) New() any {
	return &gateway.EventWebhooksUpdate{}
}

func (h *gatewayHandlerWebhooksUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventWebhooksUpdate)

	client.EventManager().DispatchEvent(&events.WebhooksUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildId:      payload.GuildID,
		ChannelID:    payload.ChannelID,
	})
}
