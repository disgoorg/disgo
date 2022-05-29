package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerInviteDelete struct{}

func (h *gatewayHandlerInviteDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeInviteDelete
}

func (h *gatewayHandlerInviteDelete) New() any {
	return &discord.GatewayEventInviteDelete{}
}

func (h *gatewayHandlerInviteDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*discord.GatewayEventInviteDelete)

	client.EventManager().DispatchEvent(&events.InviteDelete{
		GenericInvite: &events.GenericInvite{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildID:      payload.GuildID,
			ChannelID:    payload.ChannelID,
			Code:         payload.Code,
		},
	})
}
