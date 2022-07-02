package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerTypingStart struct {}

func (h *gatewayHandlerTypingStart) EventType() gateway.EventType {
	return gateway.EventTypeTypingStart
}

func (h *gatewayHandlerTypingStart) New() any {
	return &gateway.EventTypingStart{}
}

func (h *gatewayHandlerTypingStart) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventTypingStart)

	client.EventManager().DispatchEvent(&events.UserTypingStart{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		ChannelID:    payload.ChannelID,
		GuildID:      payload.GuildID,
		UserID:       payload.UserID,
		Timestamp:    payload.Timestamp,
	})

	if payload.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMUserTypingStart{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    payload.ChannelID,
			UserID:       payload.UserID,
			Timestamp:    payload.Timestamp,
		})
	} else {
		client.Caches().Members().Put(*payload.GuildID, payload.UserID, *payload.Member)
		client.EventManager().DispatchEvent(&events.GuildMemberTypingStart{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    payload.ChannelID,
			UserID:       payload.UserID,
			GuildID:      *payload.GuildID,
			Timestamp:    payload.Timestamp,
			Member:       *payload.Member,
		})
	}
}
