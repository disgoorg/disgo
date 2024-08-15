package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildSoundboardSoundCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundCreate) {
	client.Caches().AddGuildSoundboardSound(event.SoundboardSound)

	client.EventManager().DispatchEvent(&events.GuildSoundboardSoundCreate{
		GenericGuildSoundboardSound: &events.GenericGuildSoundboardSound{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			SoundboardSound: event.SoundboardSound,
		},
	})
}

func gatewayHandlerGuildSoundboardSoundUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundUpdate) {
	oldSound, _ := client.Caches().GuildSoundboardSound(*event.GuildID, event.SoundID)
	client.Caches().AddGuildSoundboardSound(event.SoundboardSound)

	client.EventManager().DispatchEvent(&events.GuildSoundboardSoundUpdate{
		GenericGuildSoundboardSound: &events.GenericGuildSoundboardSound{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			SoundboardSound: event.SoundboardSound,
		},
		OldGuildSoundboardSound: oldSound,
	})
}

func gatewayHandlerGuildSoundboardSoundDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundDelete) {
	client.Caches().RemoveGuildSoundboardSound(event.GuildID, event.SoundID)

	client.EventManager().DispatchEvent(&events.GuildSoundboardSoundDelete{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		SoundID:      event.SoundID,
		GuildID:      event.GuildID,
	})
}

func gatewayHandlerGuildSoundboardSoundsUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundsUpdate) {
	for _, sound := range event {
		client.Caches().AddGuildSoundboardSound(sound)
	}

	client.EventManager().DispatchEvent(&events.GuildSoundboardSoundsUpdate{
		GenericEvent:     events.NewGenericEvent(client, sequenceNumber, shardID),
		SoundboardSounds: event,
	})
}

func gatewayHandlerSoundboardSounds(client bot.Client, sequenceNumber int, shardID int, event gateway.EventSoundboardSounds) {
	client.EventManager().DispatchEvent(&events.SoundboardSounds{
		GenericEvent:     events.NewGenericEvent(client, sequenceNumber, shardID),
		SoundboardSounds: event.SoundboardSounds,
		GuildID:          event.GuildID,
	})
}
