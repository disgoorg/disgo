package handlers

import (
	"github.com/snekROmonoro/disgo/bot"
	"github.com/snekROmonoro/disgo/events"
	"github.com/snekROmonoro/disgo/gateway"
)

func gatewayHandlerGuildBanAdd(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildBanAdd) {
	client.EventManager().DispatchEvent(&events.GuildBan{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
		User:         event.User,
	})
}

func gatewayHandlerGuildBanRemove(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildBanRemove) {
	client.EventManager().DispatchEvent(&events.GuildUnban{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      event.GuildID,
		User:         event.User,
	})
}
