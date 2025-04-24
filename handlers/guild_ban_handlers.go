package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildBanAdd(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildBanAdd) {
	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
	}

	client.EventManager.DispatchEvent(&events.GuildBan{
		GenericGuild: genericGuildEvent,
		User:         event.User,
	})
}

func gatewayHandlerGuildBanRemove(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildBanRemove) {
	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
	}

	client.EventManager.DispatchEvent(&events.GuildUnban{
		GenericGuild: genericGuildEvent,
		User:         event.User,
	})
}
