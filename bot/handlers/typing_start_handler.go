package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerTypingStart(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventTypingStart) {
	client.EventManager.DispatchEvent(&events.UserTypingStart{
		Event:        events.NewEvent(client),
		GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
		ChannelID:    event.ChannelID,
		GuildID:      event.GuildID,
		UserID:       event.UserID,
		Timestamp:    event.Timestamp,
	})

	if event.GuildID == nil {
		client.EventManager.DispatchEvent(&events.DMUserTypingStart{
			Event:        events.NewEvent(client),
			GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
			ChannelID:    event.ChannelID,
			UserID:       event.UserID,
			Timestamp:    event.Timestamp,
		})
	} else {
		var member discord.Member
		if event.Member != nil {
			member = *event.Member
		}
		client.EventManager.DispatchEvent(&events.GuildMemberTypingStart{
			Event:        events.NewEvent(client),
			GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
			ChannelID:    event.ChannelID,
			UserID:       event.UserID,
			GuildID:      *event.GuildID,
			Timestamp:    event.Timestamp,
			Member:       member,
		})
	}
}
