package handlers

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerRaw(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventRaw) {
	client.EventManager.DispatchEvent(&events.Raw{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		EventRaw:     event,
	})
}

func gatewayHandlerHeartbeatAck(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventHeartbeatAck) {
	client.EventManager.DispatchEvent(&events.HeartbeatAck{
		GenericEvent:      events.NewGenericEvent(client, sequenceNumber, shardID),
		EventHeartbeatAck: event,
	})
}

func gatewayHandlerReady(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventReady) {
	if err := client.Caches.SetSelfUser(event.User); err != nil {
		client.Logger.Error("failed to set self user to cache", slog.Any("err", err), slog.String("user_id", event.User.ID.String()))
	}

	for _, guild := range event.Guilds {
		if err := client.Caches.SetGuildUnready(guild.ID, true); err != nil {
			client.Logger.Error("failed to set guild unready to cache", slog.Any("err", err), slog.String("guild_id", guild.ID.String()))
		}
	}

	client.EventManager.DispatchEvent(&events.Ready{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		EventReady:   event,
	})
}

func gatewayHandlerResumed(client *bot.Client, sequenceNumber int, shardID int, _ gateway.EventData) {
	client.EventManager.DispatchEvent(&events.Resumed{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
	})
}
