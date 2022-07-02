package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func gatewayHandlerVoiceStateUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventVoiceStateUpdate) {
	member := event.Member

	oldVoiceState, oldOk := client.Caches().VoiceStates().Get(event.GuildID, event.UserID)
	if event.ChannelID == nil {
		client.Caches().VoiceStates().Remove(event.GuildID, event.UserID)
	} else {
		client.Caches().VoiceStates().Put(event.GuildID, event.UserID, event.VoiceState)
	}
	client.Caches().Members().Put(event.GuildID, event.UserID, member)

	genericGuildVoiceEvent := &events.GenericGuildVoiceState{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		VoiceState:   event.VoiceState,
		Member:       member,
	}

	client.EventManager().DispatchEvent(&events.GuildVoiceStateUpdate{
		GenericGuildVoiceState: genericGuildVoiceEvent,
		OldVoiceState:          oldVoiceState,
	})

	if oldOk && oldVoiceState.ChannelID != nil && event.ChannelID != nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceMove{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else if (oldOk || oldVoiceState.ChannelID == nil) && event.ChannelID != nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceJoin{
			GenericGuildVoiceState: genericGuildVoiceEvent,
		})
	} else if event.ChannelID == nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceLeave{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else {
		client.Logger().Warnf("could not decide which GuildVoice to fire")
	}
}

func gatewayHandlerVoiceServerUpdate(client bot.Client, sequenceNumber int, shardID int, event gateway.EventVoiceServerUpdate) {
	client.EventManager().DispatchEvent(&events.VoiceServerUpdate{
		GenericEvent:           events.NewGenericEvent(client, sequenceNumber, shardID),
		EventVoiceServerUpdate: event,
	})
}
