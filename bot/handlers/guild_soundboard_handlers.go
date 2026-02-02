package handlers

import (
	"errors"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildSoundboardSoundCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundCreate) {
	if err := client.Caches.AddGuildSoundboardSound(event.SoundboardSound); err != nil {
		client.Logger.Error("failed to add guild soundboard sound to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("sound_id", event.SoundboardSound.SoundID.String()))
	}

	client.EventManager.DispatchEvent(&events.GuildSoundboardSoundCreate{
		GenericGuildSoundboardSound: &events.GenericGuildSoundboardSound{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			SoundboardSound: event.SoundboardSound,
		},
	})
}

func gatewayHandlerGuildSoundboardSoundUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundUpdate) {
	oldSound, err := client.Caches.GuildSoundboardSound(*event.GuildID, event.SoundID)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		client.Logger.Error("failed to get soundboard sound from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("sound_id", event.SoundID.String()))
	}
	if err := client.Caches.AddGuildSoundboardSound(event.SoundboardSound); err != nil {
		client.Logger.Error("failed to add guild soundboard sound to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("sound_id", event.SoundboardSound.SoundID.String()))
	}

	client.EventManager.DispatchEvent(&events.GuildSoundboardSoundUpdate{
		GenericGuildSoundboardSound: &events.GenericGuildSoundboardSound{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			SoundboardSound: event.SoundboardSound,
		},
		OldGuildSoundboardSound: oldSound,
	})
}

func gatewayHandlerGuildSoundboardSoundDelete(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundDelete) {
	if _, err := client.Caches.RemoveGuildSoundboardSound(event.GuildID, event.SoundID); err != nil {
		client.Logger.Error("failed to remove guild soundboard sound from cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("sound_id", event.SoundID.String()))
	}

	client.EventManager.DispatchEvent(&events.GuildSoundboardSoundDelete{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		SoundID:      event.SoundID,
		GuildID:      event.GuildID,
	})
}

func gatewayHandlerGuildSoundboardSoundsUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundsUpdate) {
	for _, sound := range event.SoundboardSounds {
		if err := client.Caches.AddGuildSoundboardSound(sound); err != nil {
			client.Logger.Error("failed to add guild soundboard sound to cache", slog.Any("err", err), slog.String("guild_id", event.GuildID.String()), slog.String("sound_id", sound.SoundID.String()))
		}
	}

	client.EventManager.DispatchEvent(&events.GuildSoundboardSoundsUpdate{
		GenericEvent:     events.NewGenericEvent(client, sequenceNumber, shardID),
		SoundboardSounds: event.SoundboardSounds,
		GuildID:          event.GuildID,
	})
}

func gatewayHandlerSoundboardSounds(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventSoundboardSounds) {
	client.EventManager.DispatchEvent(&events.SoundboardSounds{
		GenericEvent:     events.NewGenericEvent(client, sequenceNumber, shardID),
		SoundboardSounds: event.SoundboardSounds,
		GuildID:          event.GuildID,
	})
}
