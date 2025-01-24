package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerInviteCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventInviteCreate) {
	client.EventManager().DispatchEvent(&events.InviteCreate{
		GenericEvent:      events.NewGenericEvent(client, sequenceNumber, shardID),
		EventInviteCreate: event,
	})
}

func gatewayHandlerInviteDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventInviteDelete) {
	client.EventManager().DispatchEvent(&events.InviteDelete{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
		ChannelID:    event.ChannelID,
		Code:         event.Code,
	})
}
