package handlers

import (
	"errors"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildScheduledEventCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildScheduledEventCreate) {
	if err := client.Caches.AddGuildScheduledEvent(event.GuildScheduledEvent); err != nil {
		client.Logger.Error("failed to add guild scheduled event to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("event_id", event.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.GuildScheduledEventCreate{
		GenericGuildScheduledEvent: &events.GenericGuildScheduledEvent{
			GenericEvent:   events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduled: event.GuildScheduledEvent,
		},
	})
}

func gatewayHandlerGuildScheduledEventUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildScheduledEventUpdate) {
	oldGuildScheduledEvent, err := client.Caches.GuildScheduledEvent(event.GuildID, event.ID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get guild scheduled event from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("event_id", event.ID.String()))
	}
	if err := client.Caches.AddGuildScheduledEvent(event.GuildScheduledEvent); err != nil {
		client.Logger.Error("failed to add guild scheduled event to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("event_id", event.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.GuildScheduledEventUpdate{
		GenericGuildScheduledEvent: &events.GenericGuildScheduledEvent{
			GenericEvent:   events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduled: event.GuildScheduledEvent,
		},
		OldGuildScheduled: oldGuildScheduledEvent,
	})
}

func gatewayHandlerGuildScheduledEventDelete(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildScheduledEventDelete) {
	if _, err := client.Caches.RemoveGuildScheduledEvent(event.GuildID, event.ID); err != nil {
		client.Logger.Error("failed to remove guild scheduled event from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("event_id", event.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.GuildScheduledEventDelete{
		GenericGuildScheduledEvent: &events.GenericGuildScheduledEvent{
			GenericEvent:   events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduled: event.GuildScheduledEvent,
		},
	})
}

func gatewayHandlerGuildScheduledEventUserAdd(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildScheduledEventUserAdd) {
	client.EventManager.DispatchEvent(&events.GuildScheduledEventUserAdd{
		GenericGuildScheduledEventUser: &events.GenericGuildScheduledEventUser{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduledEventID: event.GuildScheduledEventID,
			UserID:                event.UserID,
			GuildID:               event.GuildID,
		},
	})
}

func gatewayHandlerGuildScheduledEventUserRemove(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildScheduledEventUserRemove) {
	client.EventManager.DispatchEvent(&events.GuildScheduledEventUserRemove{
		GenericGuildScheduledEventUser: &events.GenericGuildScheduledEventUser{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduledEventID: event.GuildScheduledEventID,
			UserID:                event.UserID,
			GuildID:               event.GuildID,
		},
	})
}
