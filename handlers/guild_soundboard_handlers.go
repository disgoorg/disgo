package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerGuildSoundboardSoundCreate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundCreate) {
	client.EventManager().DispatchEvent(&events.GuildSoundboardSoundCreate{
		GenericGuildSoundboardSound: &events.GenericGuildSoundboardSound{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			SoundboardSound: event.SoundboardSound,
		},
	})
}

func gatewayHandlerGuildSoundboardSoundUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundUpdate) {
	client.EventManager().DispatchEvent(&events.GuildSoundboardSoundUpdate{
		GenericGuildSoundboardSound: &events.GenericGuildSoundboardSound{
			GenericEvent:    events.NewGenericEvent(client, sequenceNumber, shardID),
			SoundboardSound: event.SoundboardSound,
		},
	})
}

func gatewayHandlerGuildSoundboardSoundDelete(client bot.Client, sequenceNumber int, shardID int, event gateway.EventGuildSoundboardSoundDelete) {
	client.EventManager().DispatchEvent(&events.GuildSoundboardSoundDelete{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		SoundID:      event.SoundID,
		GuildID:      event.GuildID,
	})
}

func gatewayHandlerSoundboardSounds(client bot.Client, sequenceNumber int, shardID int, event gateway.EventSoundboardSounds) {
	client.EventManager().DispatchEvent(&events.SoundboardSounds{
		GenericEvent:     events.NewGenericEvent(client, sequenceNumber, shardID),
		SoundboardSounds: event.SoundboardSounds,
		GuildID:          event.GuildID,
	})
}
