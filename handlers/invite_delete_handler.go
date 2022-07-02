package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerInviteDelete struct{}

func (h *gatewayHandlerInviteDelete) EventType() gateway.EventType {
	return gateway.EventTypeInviteDelete
}

func (h *gatewayHandlerInviteDelete) New() any {
	return &gateway.EventInviteDelete{}
}

func (h *gatewayHandlerInviteDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventInviteDelete)

	client.EventManager().DispatchEvent(&events.InviteDelete{
		GenericInvite: &events.GenericInvite{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			ChannelID:    payload.ChannelID,
			Code:         payload.Code,
		},
	})
}
