package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildSoundboardSoundCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundCreate) {
	client.Caches.AddGuildSoundboardSound(event.SoundboardSound)

	client.EventManager.DispatchEvent(&events.GuildSoundboardSoundCreate{
		GenericGuildSoundboardSound: &events.GenericGuildSoundboardSound{
			Event:           events.NewEvent(client),
			GatewayEvent:    events.NewGatewayEvent(sequenceNumber, shardID),
			SoundboardSound: event.SoundboardSound,
		},
	})
}

func gatewayHandlerGuildSoundboardSoundUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundUpdate) {
	oldSound, _ := client.Caches.GuildSoundboardSound(*event.GuildID, event.SoundID)
	client.Caches.AddGuildSoundboardSound(event.SoundboardSound)

	client.EventManager.DispatchEvent(&events.GuildSoundboardSoundUpdate{
		GenericGuildSoundboardSound: &events.GenericGuildSoundboardSound{
			Event:           events.NewEvent(client),
			GatewayEvent:    events.NewGatewayEvent(sequenceNumber, shardID),
			SoundboardSound: event.SoundboardSound,
		},
		OldGuildSoundboardSound: oldSound,
	})
}

func gatewayHandlerGuildSoundboardSoundDelete(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundDelete) {
	client.Caches.RemoveGuildSoundboardSound(event.GuildID, event.SoundID)

	client.EventManager.DispatchEvent(&events.GuildSoundboardSoundDelete{
		Event:        events.NewEvent(client),
		GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
		SoundID:      event.SoundID,
		GuildID:      event.GuildID,
	})
}

func gatewayHandlerGuildSoundboardSoundsUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundsUpdate) {
	for _, sound := range event.SoundboardSounds {
		client.Caches.AddGuildSoundboardSound(sound)
	}

	client.EventManager.DispatchEvent(&events.GuildSoundboardSoundsUpdate{
		Event:            events.NewEvent(client),
		GatewayEvent:     events.NewGatewayEvent(sequenceNumber, shardID),
		SoundboardSounds: event.SoundboardSounds,
		GuildID:          event.GuildID,
	})
}

func gatewayHandlerSoundboardSounds(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventSoundboardSounds) {
	client.EventManager.DispatchEvent(&events.SoundboardSounds{
		Event:            events.NewEvent(client),
		GatewayEvent:     events.NewGatewayEvent(sequenceNumber, shardID),
		SoundboardSounds: event.SoundboardSounds,
		GuildID:          event.GuildID,
	})
}
