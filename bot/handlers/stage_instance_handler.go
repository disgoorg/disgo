package handlers

import (
	"errors"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerStageInstanceCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventStageInstanceCreate) {
	if err := client.Caches.AddStageInstance(event.StageInstance); err != nil {
		client.Logger.Error("failed to add stage instance to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("stage_instance_id", event.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.StageInstanceCreate{
		GenericStageInstance: &events.GenericStageInstance{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			StageInstanceID: event.ID,
			StageInstance:   event.StageInstance,
		},
	})
}

func gatewayHandlerStageInstanceUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventStageInstanceUpdate) {
	oldStageInstance, err := client.Caches.StageInstance(event.GuildID, event.ID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get stage instance from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("stage_instance_id", event.ID.String()))
	}
	if err := client.Caches.AddStageInstance(event.StageInstance); err != nil {
		client.Logger.Error("failed to add stage instance to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("stage_instance_id", event.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.StageInstanceUpdate{
		GenericStageInstance: &events.GenericStageInstance{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			StageInstanceID: event.ID,
			StageInstance:   event.StageInstance,
		},
		OldStageInstance: oldStageInstance,
	})
}

func gatewayHandlerStageInstanceDelete(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventStageInstanceDelete) {
	if _, err := client.Caches.RemoveStageInstance(event.GuildID, event.ID); err != nil {
		client.Logger.Error("failed to remove stage instance from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("stage_instance_id", event.ID.String()))
	}

	client.EventManager.DispatchEvent(&events.StageInstanceDelete{
		GenericStageInstance: &events.GenericStageInstance{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			StageInstanceID: event.ID,
			StageInstance:   event.StageInstance,
		},
	})
}
