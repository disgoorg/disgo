package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerGuildIntegrationsUpdate struct{}

func (h *gatewayHandlerGuildIntegrationsUpdate) EventType() gateway.EventType {
	return gateway.EventTypeGuildIntegrationsUpdate
}

func (h *gatewayHandlerGuildIntegrationsUpdate) New() any {
	return &gateway.EventGuildIntegrationsUpdate{}
}

func (h *gatewayHandlerGuildIntegrationsUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventGuildIntegrationsUpdate)

	client.EventManager().DispatchEvent(&events.GuildIntegrationsUpdate{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      payload.GuildID,
	})
}
