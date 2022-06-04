package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type gatewayHandlerVoiceStateUpdate struct{}

func (h *gatewayHandlerVoiceStateUpdate) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeVoiceStateUpdate
}

func (h *gatewayHandlerVoiceStateUpdate) New() any {
	return &discord.VoiceStateUpdate{}
}

func (h *gatewayHandlerVoiceStateUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	voiceStateUpdate := *v.(*discord.VoiceStateUpdate)

	oldVoiceState, oldOk := client.Caches().VoiceStates().Get(voiceStateUpdate.GuildID, voiceStateUpdate.UserID)
	if voiceStateUpdate.ChannelID == nil {
		client.Caches().VoiceStates().Remove(voiceStateUpdate.GuildID, voiceStateUpdate.UserID)
	} else {
		client.Caches().VoiceStates().Put(voiceStateUpdate.GuildID, voiceStateUpdate.UserID, voiceStateUpdate.VoiceState)
	}
	client.Caches().Members().Put(voiceStateUpdate.GuildID, voiceStateUpdate.UserID, voiceStateUpdate.Member)

	if voiceStateUpdate.UserID == client.ID() && client.VoiceManager() != nil {
		client.VoiceManager().HandleVoiceStateUpdate(voiceStateUpdate)
	}

	genericGuildVoiceEvent := &events.GenericGuildVoiceState{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		VoiceState:   voiceStateUpdate.VoiceState,
		Member:       voiceStateUpdate.Member,
	}

	client.EventManager().DispatchEvent(&events.GuildVoiceStateUpdate{
		GenericGuildVoiceState: genericGuildVoiceEvent,
		OldVoiceState:          oldVoiceState,
	})

	if oldOk && oldVoiceState.ChannelID != nil && voiceStateUpdate.ChannelID != nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceMove{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else if (oldOk || oldVoiceState.ChannelID == nil) && voiceStateUpdate.ChannelID != nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceJoin{
			GenericGuildVoiceState: genericGuildVoiceEvent,
		})
	} else if voiceStateUpdate.ChannelID == nil {
		client.EventManager().DispatchEvent(&events.GuildVoiceLeave{
			GenericGuildVoiceState: genericGuildVoiceEvent,
			OldVoiceState:          oldVoiceState,
		})
	} else {
		client.Logger().Warnf("could not decide which GuildVoice to fire")
	}
}
