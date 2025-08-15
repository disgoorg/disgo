package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerVoiceChannelEffectSend(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventVoiceChannelEffectSend) {
	client.EventManager.DispatchEvent(&events.GuildVoiceChannelEffectSend{
		Event:                       events.NewEvent(client),
		GatewayEvent:                events.NewGatewayEvent(sequenceNumber, shardID),
		EventVoiceChannelEffectSend: event,
	})
}

func gatewayHandlerVoiceStateUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventVoiceStateUpdate) {
	member := event.Member

	oldVoiceState, oldOk := client.Caches.VoiceState(event.GuildID, event.UserID)
	if event.ChannelID == nil {
		client.Caches.RemoveVoiceState(event.GuildID, event.UserID)
	} else {
		client.Caches.AddVoiceState(event.VoiceState)
	}
	client.Caches.AddMember(member)

	if event.UserID == client.ID() && client.VoiceManager != nil {
		client.VoiceManager.HandleVoiceStateUpdate(event)
	}

	genericGuildVoiceEvent := &events.GenericGuildVoiceState{
		Event:        events.NewEvent(client),
		GatewayEvent: events.NewGatewayEvent(sequenceNumber, shardID),
		VoiceState:   event.VoiceState,
		Member:       member,
	}

	client.EventManager.DispatchEvent(&events.GuildVoiceStateUpdate{
		GenericGuildVoiceState: genericGuildVoiceEvent,
		OldVoiceState:          oldVoiceState,
	})

	if oldOk && oldVoiceState.ChannelID != nil && event.ChannelID != nil {
		client.EventManager.DispatchEvent(&events.GuildVoiceMove{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else if (oldOk || oldVoiceState.ChannelID == nil) && event.ChannelID != nil {
		client.EventManager.DispatchEvent(&events.GuildVoiceJoin{
			GenericGuildVoiceState: genericGuildVoiceEvent,
		})
	} else if event.ChannelID == nil {
		client.EventManager.DispatchEvent(&events.GuildVoiceLeave{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else {
		client.Logger.Warn("could not decide which GuildVoice to fire")
	}
}

func gatewayHandlerVoiceServerUpdate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventVoiceServerUpdate) {
	if client.VoiceManager != nil {
		client.VoiceManager.HandleVoiceServerUpdate(event)
	}

	client.EventManager.DispatchEvent(&events.VoiceServerUpdate{
		Event:                  events.NewEvent(client),
		GatewayEvent:           events.NewGatewayEvent(sequenceNumber, shardID),
		EventVoiceServerUpdate: event,
	})
}
