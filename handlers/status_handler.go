package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerRaw(client bot.Client, sequenceNumber int, shardID int, event gateway.EventRaw) {
	client.EventManager().DispatchEvent(&events.Raw{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		EventRaw:     event,
	})
}

func gatewayHandlerReady(client bot.Client, sequenceNumber int, shardID int, event gateway.EventReady) {
	client.Caches().PutSelfUser(event.User)

	for _, guild := range event.Guilds {
		client.Caches().Guilds().SetUnready(shardID, guild.ID)
	}

	client.EventManager().DispatchEvent(&events.Ready{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		EventReady:   event,
	})
}

func gatewayHandlerResumed(client bot.Client, sequenceNumber int, shardID int, _ gateway.EventData) {
	client.EventManager().DispatchEvent(&events.Resumed{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
	})
}
