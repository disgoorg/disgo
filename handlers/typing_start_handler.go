package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerTypingStart(client bot.Client, sequenceNumber int, shardID int, event gateway.EventTypingStart) {
	client.EventManager().DispatchEvent(&events.UserTypingStart{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		ChannelID:    event.ChannelID,
		GuildID:      event.GuildID,
		UserID:       event.UserID,
		Timestamp:    event.Timestamp,
	})

	if event.GuildID == nil {
		client.EventManager().DispatchEvent(&events.DMUserTypingStart{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    event.ChannelID,
			UserID:       event.UserID,
			Timestamp:    event.Timestamp,
		})
	} else {
		client.Caches().AddMember(*event.Member)
		client.EventManager().DispatchEvent(&events.GuildMemberTypingStart{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			ChannelID:    event.ChannelID,
			UserID:       event.UserID,
			GuildID:      *event.GuildID,
			Timestamp:    event.Timestamp,
			Member:       *event.Member,
		})
	}
}
